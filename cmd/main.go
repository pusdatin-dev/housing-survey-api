package main

import (
	"context"
	"fmt"
	"housing-survey-api/config"
	"housing-survey-api/controllers"
	appcontext "housing-survey-api/internal/context"
	"housing-survey-api/middleware"
	"housing-survey-api/models"
	"housing-survey-api/routes"
	"housing-survey-api/seed"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	db := config.InitDB(cfg)

	if cfg.DBSeed {
		fmt.Println("Running database seeder...")
		seed.RunSeeder(db, cfg)
	}
	fmt.Println("Finish database seeder")

	appCtx := &appcontext.AppContext{
		DB:     db,
		Config: cfg,
	}

	// Initialize services
	ctrl := controllers.InitControllers(appCtx)

	app := fiber.New(fiber.Config{
		DisableDefaultDate:           true,
		DisablePreParseMultipartForm: true,
	})

	// Setup middleware and routes
	middleware.InitMiddleware(appCtx)
	routes.SetupRoutes(app, ctrl)
	routes.PrintRoutes(app)

	// Graceful shutdown
	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	// Insert audit log
	insertShutdownLog(db)

	// Close DB connection
	closeDBConnection(db)

	log.Println("✅ Shutdown complete.")
	log.Println("✅ Server exited gracefully")
}

func insertShutdownLog(db *gorm.DB) {
	action := "shutdown"
	entity := "server"
	email := "housing-survey-api"
	actor := "systems"
	detail := "Graceful shutdown triggered"
	logEntry := models.AuditLog{
		Action:    &action,
		Entity:    &entity,
		Email:     &email,
		IP:        &actor,
		Detail:    &detail,
		CreatedAt: time.Now(),
	}
	if err := db.Create(&logEntry).Error; err != nil {
		log.Printf("⚠️ Failed to insert shutdown audit log: %v", err)
	} else {
		log.Println("📒 Shutdown audit log saved.")
	}
}

func closeDBConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("⚠️ Could not get sql.DB: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("⚠️ Failed to close database: %v", err)
	} else {
		log.Println("🔌 Database connection closed.")
	}
}
