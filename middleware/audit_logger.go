package middleware

import (
	"fmt"
	"housing-survey-api/utils"
	"strings"
	"time"

	"housing-survey-api/config"
	"housing-survey-api/models"

	"github.com/gofiber/fiber/v2"
)

func AuditLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Extract context data safely
		userID, _ := utils.GetUserIDFromContext(c)
		email, _ := utils.GetUserEmailFromContext(c)
		name, _ := utils.GetUserNameFromContext(c)
		role, _ := utils.GetRoleNameFromContext(c)

		ip := c.IP()
		url := c.OriginalURL()
		method := c.Method()

		// START phase log
		LogAudit("START", userID, name, email, ip, method, url, fiber.StatusOK, "")
		fmt.Printf("[API] %s %s from %s | user: %s (ID: %s, role: %s)\n", method, url, ip, email, userID, role)

		// Proceed with handler
		err := c.Next()

		// Collect status and error
		status := c.Response().StatusCode()
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}

		// Log END
		LogAudit("END", userID, name, email, ip, method, url, status, errMsg)
		duration := time.Since(start)
		fmt.Printf("[API Done] %s %s (%v)\n", method, url, duration)

		return err
	}
}

func LogAudit(phase, userID, name, email, ip, method, url string, status int, errMsg string) {
	action := fmt.Sprintf("%s_%s", phase, strings.ToUpper(method))
	entity := url

	detail := name
	if phase == "END" {
		detail = fmt.Sprintf("status: %d", status)
		if errMsg != "" {
			detail += fmt.Sprintf(", error: %s", errMsg)
		}
	}

	log := models.AuditLog{
		UserID:    pointer(userID),
		Email:     pointer(email),
		IP:        pointer(ip),
		Action:    &action,
		Entity:    &entity,
		Detail:    pointer(detail),
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&log).Error; err != nil {
		fmt.Printf("Failed to insert audit log: %v\n", err)
	}
}

func pointer(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}
