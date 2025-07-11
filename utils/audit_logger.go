package utils

import (
	"github.com/gofiber/fiber/v2"
)

func GetActorEmailOrIP(c *fiber.Ctx) string {
	if c == nil {
		return "unknown"
	}

	email, ok := c.Locals("email").(string)
	if ok && email != "" {
		return email
	}
	return c.IP()
}
