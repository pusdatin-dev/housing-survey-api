package controllers

import (
	"fmt"
	"net/http"

	"housing-survey-api/models"
	"housing-survey-api/services"
	"housing-survey-api/shared"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
)

type SurveyController struct {
	Survey services.SurveyService
}

func (c *SurveyController) GetAllSurveys(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Survey.GetAllSurveys(ctx))
}
func (c *SurveyController) GetSurveyByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return utils.ToFiberJSON(ctx, c.Survey.GetSurveyDetail(ctx, id))
}
func (c *SurveyController) CreateSurvey(ctx *fiber.Ctx) error {
	var input models.SurveyInput
	if err := ctx.BodyParser(&input); err != nil {
		fmt.Println("Error parsing request body:", err)
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}

	input.Actor = utils.GetActor(ctx)
	input.Mode = shared.Create
	res := c.Survey.CreateSurvey(ctx, input)
	return utils.ToFiberJSON(ctx, res)
}

func (c *SurveyController) UpdateSurvey(ctx *fiber.Ctx) error {
	var input models.SurveyInput
	if err := ctx.BodyParser(&input); err != nil {
		fmt.Println("Error parsing request body:", err)
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}

	input.Actor = utils.GetActor(ctx)
	input.Mode = shared.Update
	res := c.Survey.UpdateSurvey(ctx, input)
	return utils.ToFiberJSON(ctx, res)
}

func (c *SurveyController) DeleteSurvey(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return utils.ToFiberJSON(ctx, c.Survey.DeleteSurvey(ctx, id))
}

// ActionSurvey handles the approval and rejection on a survey
func (c *SurveyController) ActionSurvey(ctx *fiber.Ctx) error {
	var input models.SurveyActionInput
	if err := ctx.BodyParser(&input); err != nil {
		fmt.Println("Error parsing request body:", err)
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}
	input.Actor = utils.GetActor(ctx)
	res := c.Survey.ActionSurvey(ctx, input)
	return utils.ToFiberJSON(ctx, res)
}

func (c *SurveyController) GetSurveysByResource(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Survey.GetSurveysByResource(ctx))
}

func (c *SurveyController) GetSurveysByProgramType(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Survey.GetSurveysByProgramType(ctx))
}

func (c *SurveyController) GetSurveysByVerificationStatus(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Survey.GetSurveysByVerificationStatus(ctx))
}
