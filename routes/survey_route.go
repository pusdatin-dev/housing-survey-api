package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SurveyRoutesV1(v1 fiber.Router, ctrl *controllers.SurveyController) {
	survey := v1.Group("/surveys")

	// 🔐 Auth-required routes
	survey.Post("", middleware.SurveyorHandler(ctrl.CreateSurvey)...)
	survey.Put("", middleware.SurveyorHandler(ctrl.UpdateSurvey)...)
	survey.Delete("/:id", middleware.SurveyorHandler(ctrl.DeleteSurvey)...)
	survey.Post("/action", middleware.AuthHandler(ctrl.ActionSurvey)...)
	// --> add api for infografis balai (survey	by balai->masuk,reject, pending eselon, verif), laporan per bulan,

	// 🌐 PublicAccess routes (no auth)
	survey.Get("", middleware.PublicHandler(ctrl.GetAllSurveys)...)
	survey.Get("/:id", middleware.PublicHandler(ctrl.GetSurveyByID)...)
}
