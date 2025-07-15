package routes

import (
	"housing-survey-api/controllers"

	"github.com/gofiber/fiber/v2"
)

// Fungsi untuk register route surveyor
func SurveyorRoutes(r fiber.Router, ctrl *controllers.SurveyorController) {
	r.Get("/surveyors", ctrl.GetAllSurveyors)
	// Bisa tambah POST, PUT, DELETE di sini juga!
}
