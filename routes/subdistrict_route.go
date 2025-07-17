package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SubdistrictRoutesV1(v1 fiber.Router, ctrl *controllers.SubdistrictController) {
	subdistrict := v1.Group("/subdistrict")

	// 🔐 Auth-required routes
	subdistrict.Post("", middleware.AdminHandler(ctrl.Create)...)
	subdistrict.Put("", middleware.AdminHandler(ctrl.Update)...)
	subdistrict.Delete("/:id", middleware.AdminHandler(ctrl.Delete)...)

	// 🌐 PublicAccess routes (no auth)
	subdistrict.Get("", middleware.PublicHandler(ctrl.GetAll)...)
	subdistrict.Get("/:id", middleware.PublicHandler(ctrl.GetByID)...)
}
