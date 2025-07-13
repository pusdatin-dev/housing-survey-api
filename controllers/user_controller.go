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

func GetUsers(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, models.OkResponse(200, "Get users", nil))
}

func GetUserByID(ctx *fiber.Ctx) error {
	return ctx.SendString("Get user by id")
}

func CreateUser(ctx *fiber.Ctx) error {
	return ctx.SendString("Create user")
}

func UpdateUser(ctx *fiber.Ctx) error {
	return ctx.SendString("Update user")
}

func DeleteUser(ctx *fiber.Ctx) error {
	return ctx.SendString("Delete user")
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
