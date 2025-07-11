package context

import (
	"housing-survey-api/config"

	"gorm.io/gorm"
)

type AppContext struct {
	DB     *gorm.DB
	Config *config.Config // your .env struct
}
