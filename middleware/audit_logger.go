package middleware

import (
	"fmt"
	"strings"
	"time"

	"housing-survey-api/config"
	"housing-survey-api/models"

	"github.com/gofiber/fiber/v2"
)

func AuditLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		// Copy required values
		userID, _ := c.Locals("user_id").(string)
		email, _ := c.Locals("user_email").(string)
		name, _ := c.Locals("user_name").(string)
		role, _ := c.Locals("role_name").(string)
		fmt.Printf("User ID: %s, email: %s, name: %s\n", userID, email, name)
		ip := c.IP()
		path := c.OriginalURL()
		method := c.Method()
		url := c.OriginalURL()

		// Before API
		fmt.Printf("[API] %s %s from %s | user: %s (ID: %s, role: %s)\n", method, path, ip, email, userID, role)
		insertAuditLog(userID, name, email, ip, method, url, "START")

		err := c.Next() // Continue to handler

		// After API
		insertAuditLog(userID, name, email, ip, method, url, "END")
		fmt.Println("finish logging audit for user:", userID, "email:", email, "name:", name)
		duration := time.Since(start) // Optional timing metric
		fmt.Printf("[API Done] %s %s (%v)\n", method, path, duration)

		return err
	}
}

func insertAuditLog(userIDStr, name, email, ip, method, url, phase string) {
	action := fmt.Sprintf("%s_%s", phase, strings.ToUpper(method))
	entity := url

	log := models.AuditLog{
		UserID:    pointer(userIDStr),
		Email:     pointer(email),
		IP:        pointer(ip),
		Action:    &action,
		Entity:    &entity,
		Detail:    pointer(name),
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&log).Error; err != nil {
		fmt.Printf("Failed to insert audit log: %v\n", err)
	}
}

func pointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func InjectUserAuditFields() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("created_by", c.Locals("user_name").(string))
		c.Locals("updated_by", c.Locals("user_name").(string))
		c.Locals("deleted_by", c.Locals("user_name").(string))
		PrintLocals(c)
		return c.Next()
	}
}

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := c.Method()
		path := c.OriginalURL()
		ip := c.IP()
		email, _ := c.Locals("user_email").(string)
		userID, _ := c.Locals("user_id").(string)
		role, _ := c.Locals("role_name").(string)

		// Before request
		fmt.Printf("[API] %s %s from %s | user: %s (ID: %s, role: %s)\n", method, path, ip, email, userID, role)

		// Continue to next middleware/handler
		err := c.Next()

		// After request
		duration := time.Since(start)
		fmt.Printf("[API Done] %s %s (%v)\n", method, path, duration)

		return err
	}
}
