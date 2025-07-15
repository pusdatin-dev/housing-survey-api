package seed

import (
	"fmt"
	"log"

	"housing-survey-api/models"

	"gorm.io/gorm"
)

func SurveyorSeed(db *gorm.DB) {
	fmt.Println("Running Surveyor Seeder...")
	surveyor := models.Surveyor{
		Name:    "Surveyor A",
		BalaiID: 1, // Ganti dengan ID Balai yang sesuai
	}
	if err := db.FirstOrCreate(&surveyor, models.Surveyor{Name: surveyor.Name}).Error; err != nil {
		log.Printf("Error seeding surveyor: %v", err)
	}
	fmt.Println("Finished Surveyor Seeder...")
}
