package controllers

import (
	"housing-survey-api/models"
	"housing-survey-api/services"
	"housing-survey-api/shared"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
)

type CommentController struct {
	Comment services.CommentService
}

func (c *CommentController) GetComments(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Comment.GetAllComments(ctx))
}

func (c *CommentController) GetCommentByID(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Comment.GetCommentByID(ctx))
}

func (c *CommentController) CreatePublicComment(ctx *fiber.Ctx) error {
	var input models.CommentInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ToFiberBadRequest(ctx, "Invalid input format")
	}

	input.Actor = utils.GetActor(ctx)
	input.Mode = shared.Create
	res := c.Comment.CreatePublicComment(ctx, input)
	return utils.ToFiberJSON(ctx, res)
}

func (c *CommentController) UpdateComment(ctx *fiber.Ctx) error {
	var input models.CommentInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ToFiberBadRequest(ctx, "Invalid input format")
	}

	input.Actor = utils.GetActor(ctx)
	input.Mode = shared.Update
	res := c.Comment.UpdateComment(ctx, input)
	return utils.ToFiberJSON(ctx, res)
}

func (c *CommentController) DeleteComment(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.Comment.DeleteComment(ctx, ctx.Params("id")))
}

func (c *CommentController) ActionComment(ctx *fiber.Ctx) error {
	var input models.CommentActionInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ToFiberBadRequest(ctx, "Invalid input format")
	}

	input.Actor = utils.GetActor(ctx)
	res := c.Comment.ActionComment(ctx, input)
	return utils.ToFiberJSON(ctx, res)
}
