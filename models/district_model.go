package models

import (
	"housing-survey-api/shared"
	"time"

	"gorm.io/gorm"
)

// District Master Data
// District = Kabupaten/Kota
type District struct {
	ID         uint   `gorm:"primary_key;autoIncrement"`
	Name       string `gorm:"type:text;uniqueIndex;not null"`
	ProvinceID uint   `gorm:"index;not null"`
	Province   Province

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (d *District) UpdateFromInput(input DistrictInput) {
	d.Name = input.Name
	d.ProvinceID = input.ProvinceID
	d.UpdatedBy = input.Actor
	d.UpdatedAt = time.Now()
}

func (d *District) Update(newDistrict *District) {
	d.Name = newDistrict.Name
	d.ProvinceID = newDistrict.ProvinceID
	d.UpdatedBy = newDistrict.UpdatedBy
	d.UpdatedAt = time.Now()
}

func (d *District) ToResponse() DistrictResponse {
	return DistrictResponse{
		ID:           d.ID,
		Name:         d.Name,
		ProvinceID:   d.ProvinceID,
		ProvinceName: d.Province.Name,
	}
}

func ToDistrictResponses(districts []District) []DistrictResponse {
	res := make([]DistrictResponse, len(districts))
	for i, d := range districts {
		res[i] = d.ToResponse()
	}
	return res
}

type DistrictResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ProvinceID   uint   `json:"province_id"`
	ProvinceName string `json:"province_name"`
}

type DistrictInput struct {
	ID         uint   `json:"id"`
	Name       string `json:"name" validate:"required"`
	ProvinceID uint   `json:"province_id" validate:"required"`
	Actor      string `json:"-"`
	Mode       string `json:"-"` // "create" or "update"
}

func (d *DistrictInput) Validate() error {
	customMessages := map[string]string{
		"Name.required":       "District name is required",
		"ProvinceID.required": "Province is required",
	}
	return shared.CustomValidate(d, customMessages)
}

func (d *DistrictInput) ToDistrict() District {
	now := time.Now()
	return District{
		ID:         d.ID,
		Name:       d.Name,
		ProvinceID: d.ProvinceID,
		CreatedBy:  d.Actor,
		UpdatedBy:  d.Actor,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
