package models

import (
	"log"

	"gorm.io/gorm"
)

// MigrateAll performs AutoMigrate on all models in a safe order.
func MigrateAll(db *gorm.DB) error {
	log.Println("🔄 Starting database migration...")

	err := db.Transaction(func(tx *gorm.DB) error {
		// Migrate in dependency-safe order (tables with no foreign keys first)
		if err := tx.AutoMigrate(
			&Role{},
			&ProgramType{},
			&Resource{},
			&Balai{},
			&User{},
			&Profile{},
			&Survey{},
			&Comment{},
			&AuditLog{},
		); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Printf("❌ Failed to auto-migrate models: %v", err)
		return err
	}

	log.Println("✅ Database migration complete.")
	return nil
}
