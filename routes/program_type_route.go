package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProgramTypeRoutesV1(v1 fiber.Router, ctrl *controllers.ProgramTypeController) {
	programType := v1.Group("/program_type")

	// 🔐 Auth-required routes
	programType.Post("", middleware.AdminHandler(ctrl.Create)...)
	programType.Put("", middleware.AdminHandler(ctrl.Update)...)
	programType.Delete("/:id", middleware.AdminHandler(ctrl.Delete)...)

	// 🌐 PublicAccess routes (no auth)
	programType.Get("", middleware.PublicHandler(ctrl.GetAll)...)
	programType.Get("/:id", middleware.PublicHandler(ctrl.GetByID)...)
}
