package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func RoleRoutesV1(v1 fiber.Router, ctrl *controllers.RoleController) {
	role := v1.Group("/role")

	// 🔐 Auth-required routes
	role.Post("", middleware.AdminHandler(ctrl.Create)...)
	role.Put("", middleware.AdminHandler(ctrl.Update)...)
	role.Delete("/:id", middleware.AdminHandler(ctrl.Delete)...)

	// 🌐 PublicAccess routes (no auth)
	role.Get("", middleware.PublicHandler(ctrl.GetAll)...)
	role.Get("/:id", middleware.PublicHandler(ctrl.GetByID)...)
}
