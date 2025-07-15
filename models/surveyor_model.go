package models

import (
	"time"

	"gorm.io/gorm"
)

// Struct Surveyor untuk representasi data surveyor di DB
type Surveyor struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `json:"name"`
	BalaiID uint   `json:"balai_id"`

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// Tambah field lain sesuai tabelmu
}
