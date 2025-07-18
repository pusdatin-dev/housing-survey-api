package seed

import (
	"fmt"
	"log"

	"housing-survey-api/models"

	"gorm.io/gorm"
)

func ResourceSeed(db *gorm.DB) {
	fmt.Println("Running Resource Seeder...")

	resource := []models.Resource{
		{Name: "Negara", ProgramTypeID: 1},
		{Name: "Pengembang", ProgramTypeID: 1},
		{Name: "Swadaya", ProgramTypeID: 1},
		{Name: "Gotong Royong", ProgramTypeID: 1},
		{Name: "Negara", ProgramTypeID: 2},
		{Name: "Pembiayaan", ProgramTypeID: 2},
		{Name: "Swadaya", ProgramTypeID: 2},
		{Name: "Gotong Royong", ProgramTypeID: 2},
		{Name: "Pengembang", ProgramTypeID: 3},
		{Name: "Investasi", ProgramTypeID: 3},
		{Name: "Gotong Royong", ProgramTypeID: 3},
	}
	for _, p := range resource {
		if err := db.FirstOrCreate(&p, models.Resource{Name: p.Name, ProgramTypeID: p.ProgramTypeID}).Error; err != nil {
			log.Printf("Error seeding Resource: %v", err)
		}
	}
	fmt.Println("Finished Resource Seeder...")
}
