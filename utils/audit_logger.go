package utils

import (
	"fmt"
	"strings"
	"time"

	"housing-survey-api/config"
	"housing-survey-api/models"

	"github.com/gofiber/fiber/v2"
)

const (
	RequestIDKey contextKey = "requestid"
)

func GetActor(c *fiber.Ctx) string {
	if c == nil {
		return "unknown"
	}
	if userID, err := GetUserIDFromContext(c); err == nil {
		return fmt.Sprint(userID)
	}
	if userEmail, err := GetUserEmailFromContext(c); err == nil {
		return userEmail
	}
	if userName, err := GetUserNameFromContext(c); err == nil {
		return userName
	}
	return c.IP()
}

func LogAudit(c *fiber.Ctx, action, message string) {
	userID, _ := GetUserIDFromContext(c)
	email, _ := GetUserEmailFromContext(c)
	role, _ := GetRoleNameFromContext(c)
	ip := c.IP()
	method := c.Method()
	url := c.OriginalURL()
	requestID := c.Get("X-Request-ID")

	log := models.AuditLog{
		RequestID: StringPtr(requestID),
		UserID:    StringPtr(fmt.Sprint(userID)),
		Email:     StringPtr(email),
		Role:      StringPtr(role),
		IP:        StringPtr(ip),
		Action:    StringPtr(action),
		Entity:    StringPtr(fmt.Sprintf("%s %s", method, url)),
		Detail:    StringPtr(message),
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&log).Error; err != nil {
		fmt.Printf("Failed to insert audit log: %v\n", err)
	}
}

func StringPtr(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}
