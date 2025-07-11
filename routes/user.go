package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutesV1(r fiber.Router) {
	user := r.Group("/users")
	authRequired := user.Group("", middleware.AuthRequired())
	authRequired.Post("/", controllers.CreateUser)
	authRequired.Put("/:id", controllers.UpdateUser)
	authRequired.Delete("/:id", controllers.DeleteUser)

	// 🌐 Public routes (no auth)
	user.Get("/", controllers.GetUsers)
	user.Get("/:id", controllers.GetUserByID)
}
