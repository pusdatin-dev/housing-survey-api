package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name    string    `gorm:"not null"`
	UserID  uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	BalaiID uint      `gorm:"index"`
	Balai   Balai
	SKNo    string
	SKDate  time.Time
	File    string

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
