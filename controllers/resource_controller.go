package controllers

import (
	"net/http"

	"housing-survey-api/models"
	"housing-survey-api/services"
	"housing-survey-api/shared"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
)

type ResourceController struct {
	Service services.ResourceService
}

func (c *ResourceController) GetAll(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Service.GetAll(ctx))
}

func (c *ResourceController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return utils.ToFiberJSON(ctx, c.Service.GetByID(ctx, id))
}

func (c *ResourceController) Create(ctx *fiber.Ctx) error {
	var input models.ResourceInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}
	input.Mode = shared.Create
	input.Actor = utils.GetActor(ctx)

	return utils.ToFiberJSON(ctx, c.Service.Create(ctx, &input))
}

func (c *ResourceController) Update(ctx *fiber.Ctx) error {
	var input models.ResourceInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}
	input.Mode = shared.Update
	input.Actor = utils.GetActor(ctx)

	return utils.ToFiberJSON(ctx, c.Service.Update(ctx, &input))
}

func (c *ResourceController) Delete(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Service.Delete(ctx, ctx.Params("id")))
}
