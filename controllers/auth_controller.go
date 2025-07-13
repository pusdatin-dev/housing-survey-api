package controllers

import (
	"fmt"
	"os"
	"time"

	"housing-survey-api/config"
	"housing-survey-api/models"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ToFiberBadRequest(c, "Invalid input format")
	}

	var user models.User
	if err := config.DB.Preload("Role").Preload("Profile").
		Where("email = ? AND deleted_at IS NULL", input.Email).
		First(&user).Error; err != nil {
		return utils.ToFiberFailedLogin(c)
	}
	fmt.Println("user found with email:", user.Email, user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return utils.ToFiberFailedLogin(c)
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
	fmt.Println(os.Getenv("JWT_SECRET"))
	signedToken, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	user.Token = &signedToken
	if err := config.DB.Model(&user).Update("token", user.Token).Error; err != nil {
		fmt.Println("Error updating user token: ", err)
		return utils.ToFiberJSON(c, models.ErrResponse(500, "Failed to update user token"))
	}
	fmt.Println("Successfully signed in user ", user)
	return utils.ToFiberJSON(c, models.OkResponse(fiber.StatusOK, "Login Successful",
		fiber.Map{
			"token": signedToken,
			"user":  user,
		}))
}

func Logout(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		fmt.Println("Failed to extract user_id from token:", err)
		return utils.ToFiberJSON(c, models.ErrResponse(401, "Unauthorized"))
	}

	err = config.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Where("deleted_at IS NULL").
		Update("token", nil).Error
	if err != nil {
		fmt.Println("Error updating token to null:", err)
		return utils.ToFiberJSON(c, models.ErrResponse(500, "Failed to logout"))
	}

	fmt.Println("Successfully logged out user:", userID)
	return utils.ToFiberJSON(c, models.OkResponse(200, "Logged out", nil))
}
