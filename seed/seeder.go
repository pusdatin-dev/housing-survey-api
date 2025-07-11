package seed

import (
	"housing-survey-api/config"

	"gorm.io/gorm"
)

func RunSeeder(db *gorm.DB, cfg *config.Config) {
	// Call all seed functions here
	RoleSeed(db, cfg)
	BalaiSeed(db)
	UsersSeedWithProfiles(db, cfg)
	// Add more seed functions as needed
	// For example: UserSeed(db), SurveySeed(db), etc.
	// Ensure that each seed function is defined in the respective file.
}
