package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ResourceRoutesV1(v1 fiber.Router, ctrl *controllers.ResourceController) {
	resource := v1.Group("/resource")

	// 🔐 Auth-required routes
	resource.Post("", middleware.AdminHandler(ctrl.Create)...)
	resource.Put("", middleware.AdminHandler(ctrl.Update)...)
	resource.Delete("/:id", middleware.AdminHandler(ctrl.Delete)...)

	// 🌐 PublicAccess routes (no auth)
	resource.Get("", middleware.PublicHandler(ctrl.GetAll)...)
	resource.Get("/:id", middleware.PublicHandler(ctrl.GetByID)...)
}
