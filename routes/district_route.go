package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func DistrictRoutesV1(v1 fiber.Router, ctrl *controllers.DistrictController) {
	district := v1.Group("/district")

	// 🔐 Auth-required routes
	district.Post("", middleware.AdminHandler(ctrl.Create)...)
	district.Put("", middleware.AdminHandler(ctrl.Update)...)
	district.Delete("/:id", middleware.AdminHandler(ctrl.Delete)...)

	// 🌐 PublicAccess routes (no auth)
	district.Get("", middleware.PublicHandler(ctrl.GetAll)...)
	district.Get("/:id", middleware.PublicHandler(ctrl.GetByID)...)
}
