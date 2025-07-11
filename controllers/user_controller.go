package controllers

import (
	"housing-survey-api/models"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	return utils.ToFiberJSON(c, models.OkResponse(200, "Get users", nil))
}

func GetUserByID(c *fiber.Ctx) error {
	return c.SendString("Get user by id")
}

func CreateUser(c *fiber.Ctx) error {
	return c.SendString("Create user")
}

func UpdateUser(c *fiber.Ctx) error {
	return c.SendString("Update user")
}

func DeleteUser(c *fiber.Ctx) error {
	return c.SendString("Delete user")
}
