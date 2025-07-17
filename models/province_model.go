package models

import (
	"housing-survey-api/shared"
	"time"

	"gorm.io/gorm"
)

// Province Master Data
type Province struct {
	ID   uint   `gorm:"primary_key;autoIncrement"`
	Name string `gorm:"type:text;uniqueIndex;not null"`

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (p *Province) UpdateFromInput(input ProvinceInput) {
	p.Name = input.Name
	p.UpdatedBy = input.Actor
	p.UpdatedAt = time.Now()
}

func (p *Province) Update(newProvince *Province) {
	p.Name = newProvince.Name
	p.UpdatedBy = newProvince.UpdatedBy
	p.UpdatedAt = time.Now()
}

func (p *Province) ToResponse() ProvinceResponse {
	return ProvinceResponse{
		ID:   p.ID,
		Name: p.Name,
	}
}

func ToProvinceResponses(provinces []Province) []ProvinceResponse {
	res := make([]ProvinceResponse, len(provinces))
	for i, p := range provinces {
		res[i] = p.ToResponse()
	}
	return res
}

type ProvinceInput struct {
	ID    uint   `json:"id"`
	Name  string `json:"name" validate:"required"`
	Actor string `json:"-"`
	Mode  string `json:"-"`
}

type ProvinceResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (p *ProvinceInput) Validate() error {
	custom := map[string]string{
		"Name.required": "Province name is required",
	}
	return shared.CustomValidate(p, custom)
}

func (p *ProvinceInput) ToProvince() Province {
	now := time.Now()
	return Province{
		ID:        p.ID,
		Name:      p.Name,
		CreatedBy: p.Actor,
		UpdatedBy: p.Actor,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
