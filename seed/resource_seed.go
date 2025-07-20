package seed

import (
	"fmt"
	"housing-survey-api/config"
	"log"

	"housing-survey-api/models"

	"gorm.io/gorm"
)

func ResourceSeed(db *gorm.DB, cfg *config.Config) {
	fmt.Println("Running Resource Seeder...")

	resource := []models.Resource{
		{Name: "Negara", ProgramTypeID: 1, Tag: cfg.Resource.TagNegara},
		{Name: "Pengembang", ProgramTypeID: 1, Tag: cfg.Resource.TagPengembang},
		{Name: "Swadaya", ProgramTypeID: 1, Tag: cfg.Resource.TagSwadaya},
		{Name: "Gotong Royong", ProgramTypeID: 1, Tag: cfg.Resource.TagGotongRoyong},
		{Name: "Negara", ProgramTypeID: 2, Tag: cfg.Resource.TagNegara},
		{Name: "Pembiayaan", ProgramTypeID: 2, Tag: cfg.Resource.TagPengembang},
		{Name: "Swadaya", ProgramTypeID: 2, Tag: cfg.Resource.TagSwadaya},
		{Name: "Gotong Royong", ProgramTypeID: 2, Tag: cfg.Resource.TagGotongRoyong},
		{Name: "Pengembang", ProgramTypeID: 3, Tag: cfg.Resource.TagPengembang},
		{Name: "Investasi", ProgramTypeID: 3, Tag: cfg.Resource.TagPengembang},
		{Name: "Gotong Royong", ProgramTypeID: 3, Tag: cfg.Resource.TagGotongRoyong},
	}
	for _, p := range resource {
		if err := db.FirstOrCreate(&p, models.Resource{Name: p.Name, ProgramTypeID: p.ProgramTypeID}).Error; err != nil {
			log.Printf("Error seeding Resource: %v", err)
		}
	}
	fmt.Println("Finished Resource Seeder...")
}
