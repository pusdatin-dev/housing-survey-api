package controllers

import (
	"fmt"

	"housing-survey-api/models"
	"housing-survey-api/services"
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
	fmt.Println("Create public comment")
	var input models.CommentInput
	if err := ctx.BodyParser(&input); err != nil {
		fmt.Println("Error parsing request body:", err)
		return utils.ToFiberBadRequest(ctx, "Invalid input format")
	}

	input.Actor = utils.GetActor(ctx)
	res := c.Comment.CreatePublicComment(ctx, input)
	fmt.Println("Success creating public comment")
	return utils.ToFiberJSON(ctx, res)
}

func (c *CommentController) UpdateComment(ctx *fiber.Ctx) error {
	// Logic to update a comment
	id := ctx.Params("id")
	var comment map[string]interface{}
	if err := ctx.BodyParser(&comment); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	return ctx.JSON(fiber.Map{"message": "Comment updated", "id": id, "comment": comment})
}

func (c *CommentController) DeleteComment(ctx *fiber.Ctx) error {
	// Logic to delete a comment
	id := ctx.Params("id")
	return ctx.JSON(fiber.Map{"message": "Comment deleted", "id": id})
}
