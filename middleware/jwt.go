package middleware

import (
	"errors"
	"housing-survey-api/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")
		if !strings.HasPrefix(tokenStr, "Bearer ") {
			utils.LogAudit(c, "AUTH_FAIL", "Missing Bearer token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		claims, err := ExtractMapClaims(tokenStr)
		if err != nil {
			utils.LogAudit(c, "AUTH_FAIL", "Invalid JWT token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// ✅ Inject claims into user context
		ctx := utils.SetUserContext(c.Context(), claims)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func PublicAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")
		if !strings.HasPrefix(tokenStr, "Bearer ") {
			ctx := utils.SetGuestContext(c.Context(), c.IP())
			c.SetUserContext(ctx)
		} else {
			claims, err := ExtractMapClaims(tokenStr)
			if err != nil {
				utils.LogAudit(c, "AUTH_FAIL", "Invalid JWT token")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
			}
			ctx := utils.SetUserContext(c.Context(), claims)
			c.SetUserContext(ctx)
		}

		utils.LogAudit(c, "GUEST_ACCESS", "Guest")
		return c.Next()
	}
}

func ExtractMapClaims(tokenStr string) (jwt.MapClaims, error) {
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	// 🔐 Use config from appCtx, not os.Getenv
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(appCtx.Config.Token), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot extract claims")
	}
	return claims, nil
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
			utils.LogAudit(c, "FORBIDDEN", "Admin access required")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		return c.Next()
	}
}

func SurveyorOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, err := utils.GetRoleNameFromContext(c)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to extract role")
		}

		// ✅ Use role names from appCtx
		cfg := appCtx.Config
		if role != cfg.Roles.Surveyor {
			utils.LogAudit(c, "FORBIDDEN", "Surveyor access required")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Surveyor access required",
			})
		}

		return c.Next()
	}
}
