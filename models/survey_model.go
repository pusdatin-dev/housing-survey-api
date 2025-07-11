package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var _ = pq.StringArray{}

type Survey struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;index"`
	User          User
	Address       string
	Coordinate    string `gorm:"type:text"` // lat,lng string or GeoJSON
	Type          string
	Images        pq.StringArray `gorm:"type:text[]"`
	ProvinceID    uint           `gorm:"index"`
	DistrictID    uint           `gorm:"index"`
	SubdistrictID uint           `gorm:"index"`
	VillageID     uint           `gorm:"index"`

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SurveyInput struct {
	ID            uuid.UUID      `json:"id"`
	UserID        uuid.UUID      `json:"user_id"`
	Address       string         `json:"address"`
	Coordinate    string         `json:"coordinate"` // lat,lng string or GeoJSON
	Type          string         `json:"type"`
	Images        pq.StringArray `json:"images"`
	ProvinceID    uint           `json:"province_id"`
	DistrictID    uint           `json:"district_id"`
	SubdistrictID uint           `json:"subdistrict_id"`
	VillageID     uint           `json:"village_id"`
	Actor         string         `json:"-"` // CreatedBy, UpdatedBy, DeletedBy
}

func (s *SurveyInput) ToSurvey() Survey {
	return Survey{
		ID:            s.ID,
		UserID:        s.UserID,
		Address:       s.Address,
		Coordinate:    s.Coordinate,
		Type:          s.Type,
		Images:        s.Images,
		ProvinceID:    s.ProvinceID,
		DistrictID:    s.DistrictID,
		SubdistrictID: s.SubdistrictID,
		VillageID:     s.VillageID,
		CreatedBy:     s.Actor,
		CreatedAt:     time.Now(),
		UpdatedBy:     s.Actor,
		UpdatedAt:     time.Now(),
	}
}
