package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SurveyRoutesV1(v1 fiber.Router, ctrl *controllers.SurveyController) {
	survey := v1.Group("/surveys")
	auth := middleware.New().Auth().Build()
	public := middleware.New().Public().Build()

	// 🔐 Auth-required routes
	survey.Post("/", middleware.With(ctrl.CreateSurvey, auth...)...)
	survey.Put("/", middleware.With(ctrl.UpdateSurvey, auth...)...)
	survey.Delete("/:id", middleware.With(ctrl.DeleteSurvey, auth...)...)
	survey.Post("/action", middleware.With(ctrl.ActionSurvey, auth...)...)
	// --> add api for infografis balai (survey	by balai->masuk,reject, pending eselon, verif), laporan per bulan,

	// 🌐 PublicAccess routes (no auth)
	survey.Get("/", middleware.With(ctrl.GetAllSurveys, public...)...)
	survey.Get("/:id", middleware.With(ctrl.GetSurveyByID, public...)...)

	//// 🔐 Auth-required routes
	//survey.Post("", ctrl.CreateSurvey)
	//survey.Put("", ctrl.UpdateSurvey)
	//survey.Delete("/:id", ctrl.DeleteSurvey)
	//survey.Post("/action", ctrl.ActionSurvey)
	//// --> add api for infografis balai (survey	by balai->masuk,reject, pending eselon, verif), laporan per bulan,
	//
	//// 🌐 PublicAccess routes (no auth)
	//survey.Get("", ctrl.GetAllSurveys)
	//survey.Get("/:id", ctrl.GetSurveyByID)
}
