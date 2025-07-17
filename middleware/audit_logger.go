package middleware

import (
	"context"
	"fmt"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AuditLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate or extract request ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set("X-Request-ID", requestID)

		// Inject into context (used in services/utils)
		ctx := context.WithValue(c.UserContext(), utils.RequestIDKey, requestID)
		c.SetUserContext(ctx)

		// Use c directly here
		utils.LogAudit(c, "API_ENTER", "User entered API")

		err := c.Next()

		status := c.Response().StatusCode()
		if err != nil {
			utils.LogAudit(c, "API_ERROR", fmt.Sprintf("status: %d, err: %v", status, err))
		} else {
			utils.LogAudit(c, "API_EXIT", fmt.Sprintf("status: %d", status))
		}

		return err
	}
}
