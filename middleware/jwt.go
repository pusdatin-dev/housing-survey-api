package middleware

import (
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
		c.Locals("role", claims["role"])
		c.Locals("email", claims["email"])
		return c.Next()
	}
}

func NonAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only set if user_id is not already set by AuthRequired
		if c.Locals("user_id") == nil {
			c.Locals("user_id", c.IP())
		}
		if c.Locals("role") == nil {
			c.Locals("role", "guest")
		}

		return c.Next()
	}
}
