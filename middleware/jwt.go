package middleware

import (
	"housing-survey-api/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		method := c.Method()
		url := c.OriginalURL()

		tokenStr := c.Get("Authorization")
		if !strings.HasPrefix(tokenStr, "Bearer ") {
			LogAudit("AUTH_FAIL", "", "", "", ip, method, url, fiber.StatusUnauthorized, "Missing Bearer token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// 🔐 Use config from appCtx, not os.Getenv
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(appCtx.Config.Token), nil
		})
		if err != nil || !token.Valid {
			LogAudit("AUTH_FAIL", "", "", "", ip, method, url, fiber.StatusUnauthorized, "Invalid JWT token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			LogAudit("AUTH_FAIL", "", "", "", ip, method, url, fiber.StatusUnauthorized, "Invalid claims")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// ✅ Inject claims into user context
		ctx := utils.SetUserContext(c.Context(), claims)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func PublicAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := utils.SetGuestContext(c.Context(), c.IP())
		c.SetUserContext(ctx)

		LogAudit("GUEST_ACCESS", "00000000-0000-0000-0000-000000000000", "Guest", "", c.IP(), c.Method(), c.OriginalURL(), fiber.StatusOK, "")
		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, err := utils.GetRoleNameFromContext(c)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to extract role")
		}

		// ✅ Use role names from appCtx
		cfg := appCtx.Config
		if role != cfg.Roles.SuperAdmin && role != cfg.Roles.AdminEselon1 {
			LogAudit("FORBIDDEN", "", "", "", c.IP(), c.Method(), c.OriginalURL(), fiber.StatusForbidden, "Admin access required")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		return c.Next()
	}
}
