package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || (role != os.Getenv("ROLE_SUPER_ADMIN") && role != os.Getenv("ROLE_ADMIN_ESELON_1")) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admin access required"})
		}
		return c.Next()
	}
}
