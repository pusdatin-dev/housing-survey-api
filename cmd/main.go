package main

import (
    "github.com/gofiber/fiber/v2"
    "housing-survey-api/config"
    "housing-survey-api/routes"
)

func main() {
    app := fiber.New()

    config.ConnectDatabase()
    routes.SetupRoutes(app)

    app.Listen(":8080")
}
