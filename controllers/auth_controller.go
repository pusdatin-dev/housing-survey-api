package controllers

import (
	"fmt"
	"time"

	"housing-survey-api/config"
	"housing-survey-api/models"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	Db     *gorm.DB
	Config *config.Config
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ToFiberBadRequest(ctx, "Invalid input format")
	}

	var user models.User
	if err := c.Db.Preload("Role").Preload("Profile").
		Where("email = ? AND deleted_at IS NULL", input.Email).
		First(&user).Error; err != nil {
		return utils.ToFiberFailedLogin(ctx)
	}
	fmt.Println("user found with email:", user.Email)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return utils.ToFiberFailedLogin(ctx)
	}

	claims := jwt.MapClaims{
		"user_id":    user.ID.String(),
		"role_id":    user.Role.ID,
		"role_name":  user.Role.Name,
		"user_email": user.Email,
		"user_name":  user.Profile.Name,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
		"start":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(c.Config.Token))
	user.Token = &signedToken
	if err := c.Db.Model(&user).Update("token", user.Token).Error; err != nil {
		fmt.Println("Error updating user token: ", err)
		return utils.ToFiberJSON(ctx, models.ErrResponse(500, "Failed to update user token"))
	}
	fmt.Println("Successfully signed in user ", user.Email)
	return utils.ToFiberJSON(ctx, models.OkResponse(fiber.StatusOK, "Login Successful",
		fiber.Map{
			"token": signedToken,
			"user":  user,
		}))
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		fmt.Println("Failed to extract user_id from token:", err)
		return utils.ToFiberJSON(ctx, models.ErrResponse(401, "Unauthorized"))
	}

	err = c.Db.Model(&models.User{}).
		Where("id = ?", userID).
		Where("deleted_at IS NULL").
		Update("token", nil).Error
	if err != nil {
		fmt.Println("Error updating token to null:", err)
		return utils.ToFiberJSON(ctx, models.ErrResponse(500, "Failed to logout"))
	}

	fmt.Println("Successfully logged out user:", userID)
	return utils.ToFiberJSON(ctx, models.OkResponse(200, "Logged out", nil))
}
