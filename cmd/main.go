package main

import (
	"fmt"
	"housing-survey-api/config"
	"housing-survey-api/controllers"
	"housing-survey-api/internal/context"
	"housing-survey-api/middleware"
	"housing-survey-api/routes"
	"housing-survey-api/seed"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()
	db := config.InitDB(cfg)
	if cfg.DBSeed {
		fmt.Println("Running database seeder...")
		seed.RunSeeder(db, cfg)
	}
	fmt.Println("Finish database seeder")

	appCtx := &context.AppContext{
		DB:     db,
		Config: cfg,
	}

	// Init services
	ctrl := controllers.InitControllers(appCtx)
	app := fiber.New(fiber.Config{
		DisableDefaultDate:           true,
		DisablePreParseMultipartForm: true,
	})
	middleware.InitMiddleware(appCtx)
	routes.SetupRoutes(app, ctrl)
	routes.PrintRoutes(app)
	log.Fatal(app.Listen(":8080"))
}
