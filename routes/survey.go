package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SurveyRoutesV1(v1 fiber.Router, ctrl *controllers.SurveyController) {
	survey := v1.Group("/surveys")
	// 🔐 Auth-required routes
	authRequired := survey.Group("", middleware.AuthRequired())
	authRequired.Post("/", ctrl.CreateSurvey)
	authRequired.Put("/", ctrl.UpdateSurvey)
	authRequired.Delete("/:id", ctrl.DeleteSurvey)

	// 🌐 Public routes (no auth)
	survey.Get("/", ctrl.GetAllSurveys)
	survey.Get("/:id", ctrl.GetSurveyByID)
}
