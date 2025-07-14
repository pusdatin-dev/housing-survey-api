package controllers

import (
	"fmt"
	"housing-survey-api/models"
	"housing-survey-api/services"
	"housing-survey-api/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	User services.UserService
}

func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.User.GetAllUsers(ctx))
}

func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.User.GetUserByID(ctx, ctx.Params("id")))
}

func (c *UserController) SignupUser(ctx *fiber.Ctx) error {
	var input models.UserInput
	if err := ctx.BodyParser(&input); err != nil {
		fmt.Println("Error parsing request body:", err)
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}

	input.Actor = utils.GetActor(ctx)
	res := c.User.SignupUser(ctx, input)
	return utils.ToFiberJSON(ctx, res)
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var input models.UserInput
	if err := ctx.BodyParser(&input); err != nil {
		fmt.Println("Error parsing request body:", err)
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}
	input.Actor = utils.GetActor(ctx)
	res := c.User.CreateUser(ctx, input)
	return utils.ToFiberJSON(ctx, res)
}

func (c *UserController) ApproveUser(ctx *fiber.Ctx) error {
	var input models.ApprovingUser
	if err := ctx.BodyParser(&input); err != nil {
		fmt.Println("Error parsing request body:", err)
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}

	input.Actor = utils.GetActor(ctx)
	res := c.User.ApproveUser(ctx, input)
	return utils.ToFiberJSON(ctx, res)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	var input models.UserInput
	if err := ctx.BodyParser(&input); err != nil {
		fmt.Println("Error parsing request body:", err)
		return utils.ToFiberJSON(ctx, models.ErrResponse(http.StatusBadRequest, "Invalid input"))
	}

	input.Actor = utils.GetActor(ctx)
	res := c.User.UpdateUser(ctx, input)
	return utils.ToFiberJSON(ctx, res)
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.User.DeleteUser(ctx, ctx.Params("id")))
}
