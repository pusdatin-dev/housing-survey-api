package models

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Name    string `gorm:"not null"`
	UserID  uint   `gorm:"uniqueIndex"`
	BalaiID *uint  `gorm:"index"`                                          // ✅ Nullable foreign key
	Balai   *Balai `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // ✅ Proper foreign key behavior
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
