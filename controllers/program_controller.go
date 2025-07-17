package controllers

import (
	"net/http"

	"housing-survey-api/models"
	"housing-survey-api/services"
	"housing-survey-api/shared"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
)

type ProgramController struct {
	Service services.ProgramService
}

func (c *ProgramController) GetAll(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Service.GetAll(ctx))
}

func (c *ProgramController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return utils.ToFiberJSON(ctx, c.Service.GetByID(ctx, id))
}

func (c *ProgramController) Create(ctx *fiber.Ctx) error {
	var input models.ProgramInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}
	input.Mode = shared.Create
	input.Actor = utils.GetActor(ctx)

	return utils.ToFiberJSON(ctx, c.Service.Create(ctx, &input))
}

func (c *ProgramController) Update(ctx *fiber.Ctx) error {
	var input models.ProgramInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}
	input.Mode = shared.Update
	input.Actor = utils.GetActor(ctx)

	return utils.ToFiberJSON(ctx, c.Service.Update(ctx, &input))
}

func (c *ProgramController) Delete(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Service.Delete(ctx, ctx.Params("id")))
}
