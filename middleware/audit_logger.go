package middleware

import (
	"fmt"
	"strings"
	"time"

	"housing-survey-api/config"
	"housing-survey-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AuditLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		// Copy required values
		userID, _ := c.Locals("user_id").(string)
		email, _ := c.Locals("email").(string)
		fmt.Printf("User ID: %s, email: %s", userID, email)
		ip := c.IP()
		method := c.Method()
		url := c.OriginalURL()

		// Before API
		insertAuditLog(userID, email, ip, method, url, "START")

		err := c.Next() // Continue to handler

		// After API
		insertAuditLog(userID, email, ip, method, url, "END")

		_ = time.Since(start) // Optional timing metric

		return err
	}
}

func insertAuditLog(userIDStr, email, ip, method, url, phase string) {
	var uid *uuid.UUID
	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			uid = &parsed
		}
	}

	var emailPtr *string
	if email != "" {
		emailPtr = &email
	}

	action := fmt.Sprintf("%s_%s", phase, strings.ToUpper(method))
	entity := url

	log := models.AuditLog{
		UserID:    uid,
		Email:     emailPtr,
		IP:        &ip,
		Action:    &action,
		Entity:    &entity,
		CreatedAt: time.Now(),
	}

	config.DB.Create(&log)
}

func InjectUserAuditFields() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if ok && userID != "" {
			uid, err := uuid.Parse(userID)
			if err == nil {
				c.Locals("created_by", uid)
				c.Locals("updated_by", uid)
				c.Locals("deleted_by", uid)
			}
		}
		return c.Next()
	}
}
