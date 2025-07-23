package controllers

import (
	"housing-survey-api/models"
	"housing-survey-api/services"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	Service services.AuthService
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var input models.LoginInput
	if err := ctx.BodyParser(&input); err != nil {
		utils.LogAudit(ctx, "LOGIN", "Invalid input format")
		return utils.ToFiberBadRequest(ctx, "Invalid input format")
	}
	res := c.Service.Login(ctx, input)
	return utils.ToFiberJSON(ctx, res)
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	res := c.Service.Logout(ctx)
	return utils.ToFiberJSON(ctx, res)
}

func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	var input models.RefreshInput
	if err := ctx.BodyParser(&input); err != nil {
		utils.LogAudit(ctx, "REFRESH_TOKEN", "Invalid input format")
		return utils.ToFiberBadRequest(ctx, "Invalid input format")
	}
	res := c.Service.RefreshToken(ctx, input)
	return utils.ToFiberJSON(ctx, res)
}
