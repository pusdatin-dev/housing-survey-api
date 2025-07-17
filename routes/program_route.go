package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProgramRoutesV1(v1 fiber.Router, ctrl *controllers.ProgramController) {
	program := v1.Group("/program")

	// 🔐 Auth-required routes
	program.Post("", middleware.AdminHandler(ctrl.Create)...)
	program.Put("", middleware.AdminHandler(ctrl.Update)...)
	program.Delete("/:id", middleware.AdminHandler(ctrl.Delete)...)

	// 🌐 PublicAccess routes (no auth)
	program.Get("", middleware.PublicHandler(ctrl.GetAll)...)
	program.Get("/:id", middleware.PublicHandler(ctrl.GetByID)...)
}
