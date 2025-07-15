package controllers

import (
	"housing-survey-api/services"

	"github.com/gofiber/fiber/v2"
)

// Struct SurveyorController pegang SurveyorService
type SurveyorController struct {
	Surveyor services.SurveyorService
}

// Constructor inject service
func NewSurveyorController(surveyor services.SurveyorService) *SurveyorController {
	return &SurveyorController{Surveyor: surveyor}
}

// Handler GET /surveyors
func (ctrl *SurveyorController) GetAllSurveyors(c *fiber.Ctx) error {
	// Panggil service untuk ambil semua surveyor
	surveyors, err := ctrl.Surveyor.GetAllSurveyors()
	if err != nil {
		// Jika error, balikin response 500 + pesan error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Jika sukses, balikin data surveyor
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": surveyors,
	})
}
