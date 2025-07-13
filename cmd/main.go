package main

import (
	"fmt"
	"housing-survey-api/config"
	"housing-survey-api/controllers"
	"housing-survey-api/internal/context"
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
	app := fiber.New()
	//app.Use(recover)

	routes.SetupRoutes(app, ctrl)
	for _, route := range app.GetRoutes() {
		fmt.Printf("Route: %-7s %-31s | Handlers: %d\n", route.Method, route.Path, len(route.Handlers))
	}
	log.Fatal(app.Listen(":8080"))
}
