package seed

import (
	"fmt"
	"log"

	"housing-survey-api/models"

	"gorm.io/gorm"
)

func BalaiSeed(db *gorm.DB) {
	fmt.Println("Running Balai Seeder...")
	balai := models.Balai{
		Name:          "Balai A",
		ProvinceID:    1,
		DistrictID:    1,
		SubdistrictID: 1,
		VillageID:     1,
	}
	if err := db.FirstOrCreate(&balai, models.Balai{Name: balai.Name}).Error; err != nil {
		log.Printf("Error seeding balai: %v", err)
	}
	fmt.Println("Finished Balai Seeder...")
}
