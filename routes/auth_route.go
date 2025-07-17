package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(r fiber.Router, ctrl *controllers.AuthController) {
	r.Post("/login", middleware.PublicHandler(ctrl.Login)...)
	r.Post("/logout", middleware.AuthHandler(ctrl.Logout)...)
	//r.Post("/login", ctrl.Login)
	//r.Post("/logout", ctrl.Logout)
}
