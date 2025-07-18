package models

import (
	"time"

	"housing-survey-api/shared"

	"gorm.io/gorm"
)

// Subdistrict Master Data
// Subdistrict = Kecamatan
type Subdistrict struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"type:text;index;not null"`
	DistrictID uint   `gorm:"index;not null"`
	District   District

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (s *Subdistrict) UpdateFromInput(input *SubdistrictInput) {
	s.Name = input.Name
	s.DistrictID = input.DistrictID
	s.UpdatedBy = input.Actor
	s.UpdatedAt = time.Now()
}

func (s *Subdistrict) Update(newSubdistrict *Subdistrict) {
	s.Name = newSubdistrict.Name
	s.DistrictID = newSubdistrict.DistrictID
	s.UpdatedBy = newSubdistrict.UpdatedBy
	s.UpdatedAt = time.Now()
}

func (s *Subdistrict) MarkDeleted(actor string) {
	s.DeletedBy = actor
	s.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
}

func (s *Subdistrict) ToResponse() SubdistrictResponse {
	return SubdistrictResponse{
		ID:           s.ID,
		Name:         s.Name,
		DistrictID:   s.DistrictID,
		DistrictName: s.District.Name,
	}
}

func ToSubdistrictResponses(data []Subdistrict) []SubdistrictResponse {
	res := make([]SubdistrictResponse, len(data))
	for i, d := range data {
		res[i] = d.ToResponse()
	}
	return res
}

type SubdistrictResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	DistrictID   uint   `json:"district_id"`
	DistrictName string `json:"district_name"`
}

type SubdistrictInput struct {
	ID         uint   `json:"id"`
	Name       string `json:"name" validate:"required"`
	DistrictID uint   `json:"district_id" validate:"required"`
	Actor      string `json:"-"`
	Mode       string `json:"-"`
}

func (s *SubdistrictInput) Validate() error {
	custom := map[string]string{
		"Name.required":       "Subdistrict name is required",
		"DistrictID.required": "District is required",
	}
	return shared.CustomValidate(s, custom)
}

func (s *SubdistrictInput) ToModel() Subdistrict {
	now := time.Now()
	return Subdistrict{
		ID:         s.ID,
		Name:       s.Name,
		DistrictID: s.DistrictID,
		CreatedBy:  s.Actor,
		UpdatedBy:  s.Actor,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
