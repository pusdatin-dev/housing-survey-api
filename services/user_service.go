package services

import (
	"errors"
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
	CreateUser(ctx *fiber.Ctx, input models.UserInput) models.ServiceResponse
	GetAllUsers(ctx *fiber.Ctx) models.ServiceResponse
	GetUserByID(ctx *fiber.Ctx, userID string) models.ServiceResponse
	UpdateUser(ctx *fiber.Ctx, input models.UserInput) models.ServiceResponse
	DeleteUser(ctx *fiber.Ctx, userID string) models.ServiceResponse
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
	allowedRoles := map[int][]uint{
		int(superAdmin.ID):   {adminEselon1.ID, adminBalai.ID},
		int(adminEselon1.ID): {verificatorEselon1.ID},
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

func (s *userService) GetAllUsers(ctx *fiber.Ctx) models.ServiceResponse {
	// Parse query params
	emailQuery := ctx.Query("email")
	showDeleted := ctx.Query("deleted") == "true"
	showActive := ctx.Query("active") == "true"
	showUnverified := ctx.Query("unverified") == "true"

	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Get actor role name from JWT context
	actorRoleName, err := utils.GetRoleNameFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Failed to get role from context")
	}

	// Determine allowed roles based on actor role
	var allowedRoleNames []string
	switch actorRoleName {
	case s.Config.Roles.SuperAdmin:
		allowedRoleNames = []string{s.Config.Roles.AdminEselon1, s.Config.Roles.AdminBalai}
	case s.Config.Roles.AdminEselon1:
		allowedRoleNames = []string{s.Config.Roles.VerificatorEselon1, s.Config.Roles.VerificatorBalai}
	case s.Config.Roles.AdminBalai:
		allowedRoleNames = []string{s.Config.Roles.VerificatorBalai, s.Config.Roles.Surveyor}
	default:
		return models.ForbiddenResponse("You are not authorized to view users")
	}

	// Convert allowed role names to role IDs
	var allowedRoles []models.Role
	if err := s.Db.Where("name IN ?", allowedRoleNames).Find(&allowedRoles).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve allowed roles")
	}
	if len(allowedRoles) == 0 {
		return models.OkResponse(fiber.StatusOK, "No users found", models.DataListResponse{})
	}

	var allowedRoleIDs []uint
	for _, role := range allowedRoles {
		allowedRoleIDs = append(allowedRoleIDs, role.ID)
	}

	// Build query
	db := s.Db.Model(&models.User{}).Preload("Role").Preload("Profile")
	db = db.Where("role_id IN ?", allowedRoleIDs)

	if emailQuery != "" {
		db = db.Where("email ILIKE ?", "%"+emailQuery+"%")
	}
	if showDeleted {
		db = db.Unscoped().Where("deleted_at IS NOT NULL")
	} else {
		db = db.Where("deleted_at IS NULL")
	}
	if showActive {
		db = db.Where("is_active = true")
	}
	if showUnverified {
		db = db.Where("is_active = false")
	}

	// Count total
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to count users")
	}

	// Query result
	var users []models.User
	if err := db.Limit(limit).Offset(offset).Order("created_at desc").Find(&users).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to fetch users")
	}

	return models.OkResponse(fiber.StatusOK, "Users retrieved successfully", models.DataListResponse{
		Data:       models.ToUserResponses(users),
		Total:      int(total),
		Page:       page,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	})
}

func (s *userService) GetUserByID(ctx *fiber.Ctx, userID string) models.ServiceResponse {
	parsedID, err := strconv.Atoi(userID)
	if err != nil {
		return models.BadRequestResponse("Invalid user ID format")
	}

	actorID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve user ID from context")
	}
	actorRoleName, err := utils.GetRoleNameFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve role from context")
	}

	// Self-access allowed
	if actorID == parsedID {
		return s.returnUserDetail(uint(parsedID))
	}

	// Fetch target user
	user, err := s.fetchUserDetail(uint(parsedID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("User not found")
		}
		return models.InternalServerErrorResponse("Failed to retrieve user")
	}

	// Use isAllowedToModify for RBAC
	if !s.isAllowedToModify(actorRoleName, user.Role.Name) {
		return models.ForbiddenResponse("You are not authorized to view this user")
	}

	return models.OkResponse(fiber.StatusOK, "User detail retrieved successfully", user.ToResponse())
}

func (s *userService) returnUserDetail(userID uint) models.ServiceResponse {
	user, err := s.fetchUserDetail(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("User not found")
		}
		return models.InternalServerErrorResponse("Failed to retrieve user")
	}
	return models.OkResponse(fiber.StatusOK, "User detail retrieved successfully", user.ToResponse())
}

func (s *userService) fetchUserDetail(userID uint) (*models.User, error) {
	var user models.User
	err := s.Db.Preload("Role").Preload("Profile").
		Where("id = ? AND deleted_at IS NULL", userID).
		First(&user).Error
	return &user, err
}

func (s *userService) UpdateUser(ctx *fiber.Ctx, input models.UserInput) models.ServiceResponse {
	// Get actor info
	actorID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve user ID from context")
	}
	actorRole, err := utils.GetRoleNameFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve role from context")
	}

	// Allow self-update
	isSelf := actorID == input.ID

	// Fetch target user
	user, err := s.fetchUserDetail(uint(input.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("User not found")
		}
		return models.InternalServerErrorResponse("Failed to retrieve user")
	}

	// If not self, check RBAC
	if !isSelf && !s.isAllowedToModify(actorRole, user.Role.Name) {
		return models.ForbiddenResponse("You are not authorized to update this user")
	}

	// Validate input
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	// Apply updates inside transaction
	err = s.Db.Transaction(func(tx *gorm.DB) error {
		user.Email = input.Email
		user.UpdatedBy = input.Actor
		user.UpdatedAt = time.Now()

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		user.Profile.Name = input.Name
		user.Profile.BalaiID = &input.BalaiID
		user.Profile.SKNo = input.SKNo
		user.Profile.SKDate = input.SKDate
		user.Profile.File = input.File
		user.Profile.UpdatedBy = input.Actor
		user.Profile.UpdatedAt = time.Now()

		return tx.Save(&user.Profile).Error
	})

	if err != nil {
		fmt.Println("Update transaction failed:", err)
		return models.InternalServerErrorResponse("Failed to update user and profile")
	}

	return models.OkResponse(fiber.StatusOK, "User and profile updated successfully", user.ToResponse())
}

func (s *userService) DeleteUser(ctx *fiber.Ctx, userID string) models.ServiceResponse {
	parsedID, err := strconv.Atoi(userID)
	if err != nil {
		return models.BadRequestResponse("Invalid user ID format")
	}

	actorID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Failed to get user ID from context")
	}
	actorRole, err := utils.GetRoleNameFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Failed to get role from context")
	}

	// Self-deletion allowed
	if actorID == parsedID {
		return s.deleteUser(uint(parsedID), actorID)
	}

	// Fetch target user
	targetUser, err := s.fetchUserDetail(uint(parsedID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("User not found")
		}
		return models.InternalServerErrorResponse("Failed to retrieve user")
	}

	// Check if actor has permission to delete target
	if !s.isAllowedToModify(actorRole, targetUser.Role.Name) {
		return models.ForbiddenResponse("You are not authorized to delete this user")
	}

	return s.deleteUser(uint(parsedID), actorID)
}

func (s *userService) deleteUser(userID uint, actorID int) models.ServiceResponse {
	var user models.User
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("User not found")
		}
		return models.InternalServerErrorResponse("Failed to retrieve user for deletion")
	}

	user.DeletedBy = fmt.Sprint(actorID)
	user.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}

	if err := s.Db.Save(&user).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to delete user")
	}

	return models.OkResponse(fiber.StatusOK, "User deleted successfully", nil)
}

func (s *userService) isAllowedToModify(actorRole, targetRole string) bool {
	allowedRoles := map[string][]string{
		s.Config.Roles.SuperAdmin:   {s.Config.Roles.AdminEselon1, s.Config.Roles.AdminBalai},
		s.Config.Roles.AdminEselon1: {s.Config.Roles.VerificatorEselon1, s.Config.Roles.VerificatorBalai},
		s.Config.Roles.AdminBalai:   {s.Config.Roles.VerificatorBalai, s.Config.Roles.Surveyor},
	}

	if validTargets, ok := allowedRoles[actorRole]; ok {
		for _, role := range validTargets {
			if role == targetRole {
				return true
			}
		}
	}
	return false
}
