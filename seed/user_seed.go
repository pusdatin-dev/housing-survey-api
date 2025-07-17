package seed

import (
	"fmt"
	"log"
	"strings"
	"time"

	"housing-survey-api/config"
	"housing-survey-api/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// helper to generate simple email
func generateEmail(roleName string) string {
	key := roleName
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, " ", "")
	return key + "@example.com"
}

func UsersSeedWithProfiles(db *gorm.DB, cfg *config.Config) {
	fmt.Println("Running User Seeder With Profile...")
	// Ensure Balai exists
	//var count int64
	//db.Model(&models.Balai{}).Count(&count)
	//if count == 0 {
	//	BalaiSeed(db)
	//}

	//var balai models.Balai
	//db.First(&balai)

	password, _ := bcrypt.GenerateFromPassword([]byte("3jutaRUMAH$"), bcrypt.DefaultCost)

	users := []struct {
		Email string
		Name  string
		Role  string
	}{
		{"superuser@gmail.com", "Super Admin", cfg.Roles.SuperAdmin},
		//{"admin1@gmail.com", "Admin Eselon 1", cfg.Roles.AdminEselon1},
		//{"ver1@gmail.com", "Verificator Eselon 1", cfg.Roles.VerificatorEselon1},
		//{"adminbalai@gmail.com", "Admin Balai", cfg.Roles.AdminBalai},
		//{"verbalai@gmail.com", "Verificator Balai", cfg.Roles.VerificatorBalai},
		//{"surveyor@gmail.com", "Surveyor", cfg.Roles.Surveyor},
	}

	for _, u := range users {
		var role models.Role
		if err := db.First(&role, "name = ?", u.Role).Error; err != nil {
			log.Printf("Role %s not found: %v", u.Role, err)
			continue
		}

		tx := db.Begin()
		user := models.User{
			ID:        1,
			Email:     u.Email,
			Password:  string(password),
			IsActive:  true,
			RoleID:    role.ID,
			CreatedBy: "seeder",
			UpdatedBy: "seeder",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}

		if err := tx.FirstOrCreate(&user, models.User{Email: u.Email}).Error; err != nil {
			log.Printf("Failed creating user %s: %v", u.Email, err)
			tx.Rollback()
			continue
		}

		profile := models.Profile{
			UserID: user.ID,
			Name:   u.Name,
			//BalaiID: balai.ID,
			CreatedBy: "seeder",
			UpdatedBy: "seeder",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}

		if err := tx.FirstOrCreate(&profile, models.Profile{UserID: user.ID}).Error; err != nil {
			log.Printf("Failed creating profile for %s: %v", u.Email, err)
			tx.Rollback()
			continue
		}

		tx.Commit()
	}
	fmt.Println("Finished User Seeder with Profiles...")
}
