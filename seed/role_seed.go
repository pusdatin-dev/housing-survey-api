package seed

import (
	"fmt"
	"log"

	"housing-survey-api/config"
	"housing-survey-api/models"

	"gorm.io/gorm"
)

func RoleSeed(db *gorm.DB, cfg *config.Config) {
	// This function is intended to seed the database with initial roles.
	// It should create roles like "SuperAdmin", "Admin", "User", etc.
	// The actual implementation would depend on the database and ORM being used.
	// For example, using GORM, you might do something like this:
	fmt.Println("Running Role Seeder...")
	roles := []models.Role{
		{Name: cfg.Roles.SuperAdmin},
		{Name: cfg.Roles.AdminEselon1},
		{Name: cfg.Roles.VerificatorEselon1},
		{Name: cfg.Roles.AdminBalai},
		{Name: cfg.Roles.VerificatorBalai},
		{Name: cfg.Roles.Surveyor},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, role).Error; err != nil {
			log.Printf("Error seeding role %s: %v", role.Name, err)
		}
	}
	fmt.Println("Finished Role Seeder...")
}
