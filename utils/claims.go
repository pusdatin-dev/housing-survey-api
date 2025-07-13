package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// getTokenFromContext extracts the Bearer token from Authorization header
func getTokenFromContext(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("missing or invalid Authorization header")
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

// parseTokenAndGetClaims decodes token into jwt.MapClaims
func parseTokenAndGetClaims(c *fiber.Ctx) (jwt.MapClaims, error) {
	tokenString, err := getTokenFromContext(c)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// You may want to validate signing method here
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unable to parse claims")
	}

	return claims, nil
}

func GetUserIDFromContext(c *fiber.Ctx) (string, error) {
	claims, err := parseTokenAndGetClaims(c)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(claims["user_id"]), nil
}

func GetUserEmailFromContext(c *fiber.Ctx) (string, error) {
	claims, err := parseTokenAndGetClaims(c)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(claims["user_email"]), nil
}

func GetUserNameFromContext(c *fiber.Ctx) (string, error) {
	claims, err := parseTokenAndGetClaims(c)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(claims["user_name"]), nil
}

func GetRoleIDFromContext(c *fiber.Ctx) (string, error) {
	claims, err := parseTokenAndGetClaims(c)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(claims["role_id"]), nil
}

func GetRoleNameFromContext(c *fiber.Ctx) (string, error) {
	claims, err := parseTokenAndGetClaims(c)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(claims["role_name"]), nil
}
