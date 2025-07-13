package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")
		if !strings.HasPrefix(tokenStr, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["user_id"])
		c.Locals("role_id", claims["role_id"])
		c.Locals("role_name", claims["role_name"])
		c.Locals("user_email", claims["user_email"])
		c.Locals("user_name", claims["user_name"])
		PrintLocals(c)
		return c.Next()
	}
}

func PublicAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("user_id", "00000000-0000-0000-0000-000000000000")
		c.Locals("role_name", "Guest")
		c.Locals("role_id", "0")
		c.Locals("user_email", "public")
		c.Locals("user_name", c.IP())
		PrintLocals(c)

		return c.Next()
	}
}

func PrintLocals(c *fiber.Ctx) {
	fmt.Println("Locals:")
	keys := []string{"user_id", "role_id", "role_name", "user_email", "user_name", "created_by", "updated_by", "deleted_by"}
	for _, key := range keys {
		fmt.Printf("Local[%s] = %v\n", key, c.Locals(key))
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role_name").(string)
		if !ok || (role != os.Getenv("ROLE_SUPER_ADMIN") && role != os.Getenv("ROLE_ADMIN_ESELON_1")) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admin access required"})
		}
		return c.Next()
	}
}
