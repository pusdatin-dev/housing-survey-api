package routes

import (
	"housing-survey-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoutesV1(r fiber.Router) {
	user := r.Group("/users")

	//authRequired := user.Group("", middleware.New().Auth().Build()...)
	authRequired := user.Group("")
	authRequired.Post("", controllers.CreateUser)
	authRequired.Put("", controllers.UpdateUser)
	authRequired.Delete("/:id", controllers.DeleteUser)

	// 🌐 PublicAccess routes (no auth)
	//public := user.Group("", middleware.New().Public().Build()...)
	public := user.Group("")
	public.Get("", controllers.GetUsers)
	public.Get("/:id", controllers.GetUserByID)
}
