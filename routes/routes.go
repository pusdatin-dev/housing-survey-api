package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, ctrl *controllers.ControllerRegistry) {
	api := app.Group("/api", middleware.NonAuth())
	// All v1 routes — logging + audit fields
	AuthRoutes(api) // /login, /signup

	v1 := api.Group("/v1", middleware.AuditLogger(), middleware.InjectUserAuditFields())
	UserRoutesV1(v1)
	CommentRoutes(v1, ctrl.Comment)
	SurveyRoutesV1(v1, ctrl.Survey)
	AuditLogRoutes(v1, ctrl.AuditLog)

	// Health check or homepage
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Housing Survey API is running")
	})
}
