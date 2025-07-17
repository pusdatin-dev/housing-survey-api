package controllers

import (
	"net/http"

	"housing-survey-api/models"
	"housing-survey-api/services"
	"housing-survey-api/shared"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
)

type SubdistrictController struct {
	Service services.SubdistrictService
}

func (c *SubdistrictController) GetAll(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Service.GetAll(ctx))
}

func (c *SubdistrictController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return utils.ToFiberJSON(ctx, c.Service.GetByID(ctx, id))
}

func (c *SubdistrictController) Create(ctx *fiber.Ctx) error {
	var input models.SubdistrictInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}
	input.Mode = shared.Create
	input.Actor = utils.GetActor(ctx)

	return utils.ToFiberJSON(ctx, c.Service.Create(ctx, &input))
}

func (c *SubdistrictController) Update(ctx *fiber.Ctx) error {
	var input models.SubdistrictInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}
	input.Mode = shared.Update
	input.Actor = utils.GetActor(ctx)

	return utils.ToFiberJSON(ctx, c.Service.Update(ctx, &input))
}

func (c *SubdistrictController) Delete(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Service.Delete(ctx, ctx.Params("id")))
}
