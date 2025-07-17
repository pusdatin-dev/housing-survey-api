package seed

import (
	"fmt"
	"housing-survey-api/models"
	"log"

	"gorm.io/gorm"
)

func MasterDataSeed(db *gorm.DB) {
	fmt.Println("Running Master Data Seeder...")
	prv := models.Province{
		ID:   1,
		Name: "Test",
	}
	if err := db.FirstOrCreate(&prv, models.Province{Name: prv.Name}).Error; err != nil {
		log.Printf("Error seeding province: %v", err)
	}

	dst := models.District{
		ID:         1,
		Name:       "Test",
		ProvinceID: prv.ID,
	}
	if err := db.FirstOrCreate(&dst, models.District{Name: dst.Name}).Error; err != nil {
		log.Printf("Error seeding district: %v", err)
	}

	subdst := models.Subdistrict{
		ID:         1,
		Name:       "Test",
		DistrictID: dst.ID,
	}
	if err := db.FirstOrCreate(&subdst, models.Subdistrict{Name: subdst.Name}).Error; err != nil {
		log.Printf("Error seeding subdistrict: %v", err)
	}

	vil := models.Village{
		ID:            1,
		Name:          "Test",
		SubdistrictID: subdst.ID,
	}
	if err := db.FirstOrCreate(&vil, models.Village{Name: vil.Name}).Error; err != nil {
		log.Printf("Error seeding vil: %v", err)
	}

	fmt.Println("Finished Master Data Seeder...")
}
