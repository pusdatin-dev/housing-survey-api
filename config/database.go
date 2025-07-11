package config

import (
	"fmt"
	"log"

	"housing-survey-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 🔥 Enable UUID extension
	if err := database.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}
	
	// Auto-migrate models
	err = database.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Profile{},
		&models.Balai{},
		&models.Survey{},
		&models.Comment{},
		&models.AuditLog{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	DB = database
	return database
}
