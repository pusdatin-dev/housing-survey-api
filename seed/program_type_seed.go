package seed

import (
	"fmt"
	"log"

	"housing-survey-api/models"

	"gorm.io/gorm"
)

func ProgramTypeSeed(db *gorm.DB) {
	fmt.Println("Running Program Type Seeder...")

	programtype := []models.ProgramType{
		{Name: "Pembangunan Baru"},
		{Name: "Peningkatan Kualitas"},
		{Name: "Pembangunan Baru - Upaya Eksternal"},
	}

	for _, p := range programtype {
		if err := db.FirstOrCreate(&p, models.ProgramType{Name: p.Name}).Error; err != nil {
			log.Printf("Error seeding Program Type: %v", err)
		}
	}
	fmt.Println("Finished Resource Seeder...")
}
