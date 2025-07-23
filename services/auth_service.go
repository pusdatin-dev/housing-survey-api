package services

import (
	"errors"

	"housing-survey-api/config"
	appcontext "housing-survey-api/internal/context"
	"housing-survey-api/models"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx *fiber.Ctx, input models.LoginInput) models.ServiceResponse
	Logout(ctx *fiber.Ctx) models.ServiceResponse
	RefreshToken(ctx *fiber.Ctx, input models.RefreshInput) models.ServiceResponse
}

type authService struct {
	Db     *gorm.DB
	Config *config.Config
}

func NewAuthService(appCtx *appcontext.AppContext) AuthService {
	return &authService{
		Db:     appCtx.DB,
		Config: appCtx.Config,
	}
}

func (s *authService) Login(ctx *fiber.Ctx, input models.LoginInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		utils.LogAudit(ctx, "LOGIN", err.Error())
		return models.BadRequestResponse("Invalid input")
	}

	var user models.User
	err := s.Db.Preload("Role").Preload("Profile").Where("email = ? AND deleted_at IS NULL", input.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.LogAudit(ctx, "LOGIN", "User not found")
		return models.UnauthorizedResponse("Invalid email or password")
	} else if err != nil {
		utils.LogAudit(ctx, "LOGIN", err.Error())
		return models.InternalServerErrorResponse("Login failed")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.LogAudit(ctx, "LOGIN", "Invalid password")
		return models.UnauthorizedResponse("Invalid email or password")
	}

	signedToken, err := utils.GenerateJWT(user, s.Config.Token)
	if err != nil {
		utils.LogAudit(ctx, "LOGIN", err.Error())
		return models.InternalServerErrorResponse("Token generation failed")
	}

	user.Token = &signedToken
	if err = s.Db.Model(&user).Update("token", user.Token).Error; err != nil {
		utils.LogAudit(ctx, "LOGIN", err.Error())
		return models.InternalServerErrorResponse("Failed to update user token")
	}
	utils.LogAudit(ctx, "LOGIN", "Login successful")
	return models.OkResponse(fiber.StatusOK, "Login successful", fiber.Map{
		"token": signedToken,
		"user":  user.ToResponse(),
	})
}

func (s *authService) Logout(ctx *fiber.Ctx) models.ServiceResponse {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		utils.LogAudit(ctx, "LOGOUT", "Unauthorized")
		return models.UnauthorizedResponse("Unauthorized")
	}

	err = s.Db.Model(&models.User{}).Where("id = ? AND deleted_at IS NULL", userID).Update("token", nil).Error
	if err != nil {
		utils.LogAudit(ctx, "LOGOUT", err.Error())
		return models.InternalServerErrorResponse("Failed to logout")
	}

	utils.LogAudit(ctx, "LOGOUT", "Logout successful")
	return models.OkResponse(fiber.StatusOK, "Logout successful", nil)
}

func (s *authService) RefreshToken(ctx *fiber.Ctx, input models.RefreshInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse("Invalid input")
	}

	claims, err := utils.ParseJWT(input.RefreshToken, s.Config.Token)
	if err != nil {
		return models.UnauthorizedResponse("Invalid or expired token")
	}

	userID := claims["user_id"].(string)
	var user models.User
	if err := s.Db.Preload("Role").Preload("Profile").Where("id = ?", userID).First(&user).Error; err != nil {
		return models.NotFoundResponse("User not found")
	}

	newToken, err := utils.GenerateJWT(user, s.Config.Token)
	if err != nil {
		return models.InternalServerErrorResponse("Failed to generate token")
	}

	if err = s.Db.Model(&user).Update("token", newToken).Error; err != nil {
		utils.LogAudit(ctx, "REFRESH_TOKEN", err.Error())
		return models.InternalServerErrorResponse("Failed to update user token")
	}
	return models.OkResponse(fiber.StatusOK, "Token refreshed", fiber.Map{
		"token": newToken,
		"user":  user.ToResponse(),
	})
}
