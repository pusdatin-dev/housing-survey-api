package utils

import (
	"github.com/gofiber/fiber/v2"
)

func GetActor(c *fiber.Ctx) string {
	if c == nil {
		return "unknown"
	}
	if userID, err := GetUserIDFromContext(c); err == nil {
		return userID
	}
	if userEmail, err := GetUserEmailFromContext(c); err == nil {
		return userEmail
	}
	if userName, err := GetUserNameFromContext(c); err == nil {
		return userName
	}
	return c.IP()
}
