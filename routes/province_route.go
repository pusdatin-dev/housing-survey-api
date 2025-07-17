package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProvinceRoutesV1(v1 fiber.Router, ctrl *controllers.ProvinceController) {
	province := v1.Group("/province")

	// 🔐 Auth-required routes
	province.Post("", middleware.AdminHandler(ctrl.Create)...)
	province.Put("", middleware.AdminHandler(ctrl.Update)...)
	province.Delete("/:id", middleware.AdminHandler(ctrl.Delete)...)

	// 🌐 PublicAccess routes (no auth)
	province.Get("", middleware.PublicHandler(ctrl.GetAll)...)
	province.Get("/:id", middleware.PublicHandler(ctrl.GetByID)...)
}
