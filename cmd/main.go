package main

import (
	"fmt"
	"log"
	
	"housing-survey-api/config"
	"housing-survey-api/controllers"
	"housing-survey-api/internal/context"
	"housing-survey-api/routes"
	"housing-survey-api/seed"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	app := fiber.New()
	app.Use(recover.New())

	routes.SetupRoutes(app, ctrl)
	for _, route := range app.GetRoutes() {
		fmt.Printf("Route registered: %s %s\n", route.Method, route.Path)
	}
	log.Fatal(app.Listen(":8080"))
}
