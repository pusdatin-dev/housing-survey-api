package routes

import (
	"housing-survey-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(r fiber.Router) {
	//public := r.Group("", middleware.New().Public().Build()...)
	public := r.Group("")
	public.Post("/login", controllers.Login)

	//auth := r.Group("", middleware.New().Auth().Build()...)
	auth := r.Group("")
	auth.Post("/logout", controllers.Logout)
}
