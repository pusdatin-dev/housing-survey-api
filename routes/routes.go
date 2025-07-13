package routes

import (
	"housing-survey-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, ctrl *controllers.ControllerRegistry) {
	api := app.Group("/api")

	v1 := api.Group("/v1")
	// All v1 routes — logging + audit fields
	AuthRoutes(v1) // /login, /signup
	UserRoutesV1(v1)
	CommentRoutes(v1, ctrl.Comment)
	SurveyRoutesV1(v1, ctrl.Survey)
	AuditLogRoutes(v1, ctrl.AuditLog)

	//// Health check or homepage
	//app.Get("/", func(c *fiber.Ctx) error {
	//	return c.SendString("Housing Survey API is running")
	//})
}
