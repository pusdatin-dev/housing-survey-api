package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutesV1(r fiber.Router, ctrl *controllers.UserController) {
	user := r.Group("/users")
	auth := middleware.New().Auth().Build()
	public := middleware.New().Public().Build()

	user.Post("", middleware.With(ctrl.CreateUser, auth...)...)
	user.Put("", middleware.With(ctrl.UpdateUser, auth...)...)
	user.Delete("/:id", middleware.With(ctrl.DeleteUser, auth...)...)
	user.Get("", middleware.With(ctrl.GetAllUsers, auth...)...)
	user.Get("/:id", middleware.With(ctrl.GetUserByID, auth...)...)
	user.Post("/approve", middleware.With(ctrl.ApproveUser, auth...)...)

	// 🌐 PublicAccess routes (no auth)
	user.Post("/signup", middleware.With(ctrl.SignupUser, public...)...)
	//
	//user.Post("", ctrl.CreateUser)
	//user.Put("", ctrl.UpdateUser)
	//user.Delete("/:id", ctrl.DeleteUser)
	//user.Get("", ctrl.GetAllUsers)
	//user.Get("/:id", ctrl.GetUserByID)
	//user.Post("/approve", ctrl.ApproveUser)
	//
	//🌐 PublicAccess routes (no auth)
	//user.Post("/signup", ctrl.SignupUser)

}
