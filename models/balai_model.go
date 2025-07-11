package models

import (
	"time"

	"gorm.io/gorm"
)

type Balai struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	ProvinceID    uint   `gorm:"index"`
	DistrictID    uint   `gorm:"index"`
	SubdistrictID uint   `gorm:"index"`
	VillageID     uint   `gorm:"index"`

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
