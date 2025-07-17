package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func VillageRoutesV1(v1 fiber.Router, ctrl *controllers.VillageController) {
	village := v1.Group("/village")

	// 🔐 Auth-required routes
	village.Post("", middleware.AdminHandler(ctrl.Create)...)
	village.Put("", middleware.AdminHandler(ctrl.Update)...)
	village.Delete("/:id", middleware.AdminHandler(ctrl.Delete)...)

	// 🌐 PublicAccess routes (no auth)
	village.Get("", middleware.PublicHandler(ctrl.GetAll)...)
	village.Get("/:id", middleware.PublicHandler(ctrl.GetByID)...)
}
