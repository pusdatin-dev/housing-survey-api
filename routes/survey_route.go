package routes

import (
	"housing-survey-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func SurveyRoutesV1(v1 fiber.Router, ctrl *controllers.SurveyController) {
	survey := v1.Group("/surveys")
	// 🔐 Auth-required routes
	//authRequired := survey.Group("", middleware.New().Auth().Build()...)
	authRequired := survey.Group("")
	authRequired.Post("", ctrl.CreateSurvey)
	authRequired.Put("", ctrl.UpdateSurvey)
	authRequired.Delete("/:id", ctrl.DeleteSurvey)

	// 🌐 PublicAccess routes (no auth)
	//public := survey.Group("", middleware.New().Public().Build()...)
	public := survey.Group("")
	public.Get("", ctrl.GetAllSurveys)
	public.Get("/:id", ctrl.GetSurveyByID)
}
