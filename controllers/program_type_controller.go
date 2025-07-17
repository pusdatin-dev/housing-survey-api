package controllers

import (
	"net/http"

	"housing-survey-api/models"
	"housing-survey-api/services"
	"housing-survey-api/shared"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
)

type ProgramTypeController struct {
	Service services.ProgramTypeService
}

func (c *ProgramTypeController) GetAll(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Service.GetAll(ctx))
}

func (c *ProgramTypeController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return utils.ToFiberJSON(ctx, c.Service.GetByID(ctx, id))
}

func (c *ProgramTypeController) Create(ctx *fiber.Ctx) error {
	var input models.ProgramTypeInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}
	input.Mode = shared.Create
	input.Actor = utils.GetActor(ctx)

	return utils.ToFiberJSON(ctx, c.Service.Create(ctx, &input))
}

func (c *ProgramTypeController) Update(ctx *fiber.Ctx) error {
	var input models.ProgramTypeInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}
	input.Mode = shared.Update
	input.Actor = utils.GetActor(ctx)

	return utils.ToFiberJSON(ctx, c.Service.Update(ctx, &input))
}

func (c *ProgramTypeController) Delete(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Service.Delete(ctx, ctx.Params("id")))
}
