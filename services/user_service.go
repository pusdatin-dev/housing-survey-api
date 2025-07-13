package services

import (
	"fmt"
	"housing-survey-api/config"
	"housing-survey-api/internal/context"
	"housing-survey-api/models"
	"housing-survey-api/utils"
	"slices"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	SignupUser(ctx *fiber.Ctx, input models.UserInput) models.ServiceResponse
	ApproveUser(ctx *fiber.Ctx, input models.ApprovingUser) models.ServiceResponse
}

type userService struct {
	Db     *gorm.DB
	Config *config.Config
}

func NewUserService(ctx *context.AppContext) UserService {
	return &userService{
		Db:     ctx.DB,
		Config: ctx.Config,
	}
}

func (s *userService) SignupUser(ctx *fiber.Ctx, input models.UserInput) models.ServiceResponse {
	// Validate input
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	var verificatorBalai, surveyor models.Role
	if err := s.Db.Where("name = ?", s.Config.Roles.VerificatorBalai).First(&verificatorBalai).Error; err != nil {
		return models.InternalServerErrorResponse("Verificator Balai role not found")
	}
	if err := s.Db.Where("name = ?", s.Config.Roles.Surveyor).First(&surveyor).Error; err != nil {
		return models.InternalServerErrorResponse("Surveyor role not found")
	}
	// Only allow roles: VerificatorBalai (4) and Surveyor (5)
	allowedRoles := map[uint]bool{
		verificatorBalai.ID: true, // VerificatorBalai
		surveyor.ID:         true, // Surveyor
	}
	if !allowedRoles[input.RoleID] {
		return models.BadRequestResponse("You are not allowed to create users with this role")
	}

	// Check if user already exists
	var existingUser models.User
	if err := s.Db.Where("email = ? AND deleted_at IS NULL", input.Email).First(&existingUser).Error; err == nil {
		return models.BadRequestResponse("User with this email already exists")
	}

	// Hash password
	password, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(password)

	// Build user
	user := input.ToUser()

	// Create transaction for user + profile
	err := s.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		profile := user.Profile
		profile.UserID = user.ID
		if err := tx.Create(&profile).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Println("Transaction failed:", err)
		return models.InternalServerErrorResponse("Failed to create user and profile")
	}

	return models.OkResponse(fiber.StatusCreated, "User and profile created successfully", user)
}

func (s *userService) CreateUser(ctx *fiber.Ctx, input models.UserInput) models.ServiceResponse {
	// Validate input
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	// Get actor role from JWT context
	actorRoleID, err := utils.GetRoleIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Failed to get role ID from context")
	}

	// Load required roles from DB
	var (
		adminEselon1       models.Role
		adminBalai         models.Role
		verificatorEselon1 models.Role
		superAdmin         models.Role
	)
	if err := s.Db.Where("name = ?", s.Config.Roles.SuperAdmin).First(&superAdmin).Error; err != nil {
		return models.InternalServerErrorResponse("SuperAdmin role not found")
	}
	if err := s.Db.Where("name = ?", s.Config.Roles.AdminEselon1).First(&adminEselon1).Error; err != nil {
		return models.InternalServerErrorResponse("AdminEselon1 role not found")
	}
	if err := s.Db.Where("name = ?", s.Config.Roles.AdminBalai).First(&adminBalai).Error; err != nil {
		return models.InternalServerErrorResponse("AdminBalai role not found")
	}
	if err := s.Db.Where("name = ?", s.Config.Roles.VerificatorEselon1).First(&verificatorEselon1).Error; err != nil {
		return models.InternalServerErrorResponse("VerificatorEselon1 role not found")
	}

	// Validate creation permission based on actor's role
	allowedRoles := map[string][]uint{
		strconv.Itoa(int(superAdmin.ID)):   {adminEselon1.ID, adminBalai.ID},
		strconv.Itoa(int(adminEselon1.ID)): {verificatorEselon1.ID},
	}

	if !slices.Contains(allowedRoles[actorRoleID], input.RoleID) {
		return models.BadRequestResponse("You are not allowed to create users with this role")
	}

	// Check if user already exists
	var existingUser models.User
	if err := s.Db.Where("email = ? AND deleted_at IS NULL", input.Email).First(&existingUser).Error; err == nil {
		return models.BadRequestResponse("User with this email already exists")
	}

	// Hash password
	password, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(password)

	// Build user
	user := input.ToUser()
	user.IsActive = true // automatically activate

	// Transaction to save user + profile
	err = s.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		profile := user.Profile
		profile.UserID = user.ID
		if err := tx.Create(&profile).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("Transaction failed:", err)
		return models.InternalServerErrorResponse("Failed to create user and profile")
	}

	return models.OkResponse(fiber.StatusCreated, "User and profile created successfully", user)
}

func (s *userService) ApproveUser(ctx *fiber.Ctx, input models.ApprovingUser) models.ServiceResponse {
	// Validate input
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	roleName, err := utils.GetRoleNameFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Failed to get role name from context")
	}
	// Define allowed approvals
	allowedRoles := map[string][]string{
		s.Config.Roles.AdminEselon1: {s.Config.Roles.VerificatorBalai},
		s.Config.Roles.AdminBalai:   {s.Config.Roles.Surveyor},
	}

	allowed, ok := allowedRoles[roleName]
	if !ok {
		return models.ForbiddenResponse("Your role is not allowed to approve users")
	}

	// Track results
	var (
		successfullyApproved []string
		failedApprovals      []fiber.Map
	)

	for _, userID := range input.UserIDs {
		var user models.User
		if err := s.Db.Preload("Role").
			Where("id = ? AND deleted_at IS NULL AND is_active IS false", userID).
			First(&user).Error; err != nil {
			failedApprovals = append(failedApprovals, fiber.Map{
				"user_id": userID,
				"error":   "User not found",
			})
			continue
		}

		// Check if user's role is allowed to be approved
		isAllowed := false
		for _, r := range allowed {
			if user.Role.Name == r {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			failedApprovals = append(failedApprovals, fiber.Map{
				"user_id":   userID,
				"user_role": user.Role.Name,
				"error":     "You are not allowed to approve this role",
			})
			continue
		}

		// Approve the user
		err := s.Db.Model(&user).Updates(map[string]interface{}{
			"is_active":  true,
			"updated_by": input.Actor,
			"updated_at": time.Now(),
		}).Error
		if err != nil {
			failedApprovals = append(failedApprovals, fiber.Map{
				"user_id": userID,
				"error":   "Failed to approve user",
			})
			continue
		}

		successfullyApproved = append(successfullyApproved, userID)
	}

	return models.OkResponse(fiber.StatusOK, "Batch approval completed", fiber.Map{
		"approved": successfullyApproved,
		"failed":   failedApprovals,
	})
}
