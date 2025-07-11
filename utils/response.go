package utils

import (
	"housing-survey-api/models"

	"github.com/gofiber/fiber/v2"
)

func ToFiberJSON(ctx *fiber.Ctx, res models.ServiceResponse) error {
	return ctx.Status(res.Code).JSON(res)
}

func ToFiberUnauthorized(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(models.ServiceResponse{
		Status:  true,
		Code:    fiber.StatusUnauthorized,
		Message: "Unauthorized",
	})
}

func ToFiberBadRequest(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(models.ServiceResponse{
		Status:  false,
		Code:    fiber.StatusBadRequest,
		Message: message,
	})
}

func ToFiberNotFound(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusNotFound).JSON(models.ServiceResponse{
		Status:  false,
		Code:    fiber.StatusNotFound,
		Message: message,
	})
}

func ToFiberInternalServerError(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(models.ServiceResponse{
		Status:  false,
		Code:    fiber.StatusInternalServerError,
		Message: message,
	})
}

func ToFiberFailedLogin(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(models.FailedLoginResponse())
}
