package models

import (
	"time"

	"housing-survey-api/shared"

	"gorm.io/gorm"
)

// Province Master Data
type Province struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:text;index;not null"`
	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Input Struct
type ProvinceInput struct {
	ID    uint   `json:"id" validate:"required_if=Mode update"`
	Name  string `json:"name" validate:"required"`
	Actor string `json:"-"` // set in controller
	Mode  string `json:"-"` // "create" or "update"
}

// Response Struct
type ProvinceResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ======= Methods =======

func (input *ProvinceInput) Validate() error {
	custom := map[string]string{
		"ID.required_if": "Province ID is required for update",
		"Name.required":  "Province name is required",
	}
	return shared.CustomValidate(input, custom)
}

func (input *ProvinceInput) ToModel() *Province {
	now := time.Now()
	return &Province{
		ID:        input.ID,
		Name:      input.Name,
		CreatedBy: input.Actor,
		UpdatedBy: input.Actor,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (p *Province) UpdateFromInput(input *ProvinceInput) {
	p.Name = input.Name
	p.UpdatedBy = input.Actor
	p.UpdatedAt = time.Now()
}

func (p *Province) UpdateFromModel(new *Province) {
	p.Name = new.Name
	p.UpdatedBy = new.UpdatedBy
	p.UpdatedAt = time.Now()
}

func (p *Province) MarkDeleted(actor string) {
	p.DeletedBy = actor
	p.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
}

func (p *Province) ToResponse() ProvinceResponse {
	return ProvinceResponse{
		ID:   p.ID,
		Name: p.Name,
	}
}

func ToProvinceResponses(list []Province) []ProvinceResponse {
	res := make([]ProvinceResponse, len(list))
	for i, p := range list {
		res[i] = p.ToResponse()
	}
	return res
}
