package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func BalaiRoutesV1(v1 fiber.Router, ctrl *controllers.BalaiController) {
	balai := v1.Group("/balai")

	// 🔐 Auth-required routes
	balai.Post("", middleware.AdminHandler(ctrl.Create)...)
	balai.Put("", middleware.AdminHandler(ctrl.Update)...)
	balai.Delete("/:id", middleware.AdminHandler(ctrl.Delete)...)

	// 🌐 PublicAccess routes (no auth)
	balai.Get("", middleware.PublicHandler(ctrl.GetAll)...)
	balai.Get("/:id", middleware.PublicHandler(ctrl.GetByID)...)
}
