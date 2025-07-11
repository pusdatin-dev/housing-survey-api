package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(r fiber.Router) {
	v1 := r.Group("", middleware.AuditLogger(), middleware.InjectUserAuditFields())
	v1.Post("/login", controllers.Login)
	auth := v1.Group("", middleware.AuthRequired())
	auth.Post("/logout", controllers.Logout)
}
