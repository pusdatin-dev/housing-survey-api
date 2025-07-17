package models

import (
	"housing-survey-api/shared"
	"time"

	"gorm.io/gorm"
)

// ProgramType Master Data
type ProgramType struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:text;not null"`
	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ProgramTypeInput struct {
	ID    uint   `json:"id"`
	Name  string `json:"name" validate:"required"`
	Actor string `json:"-"`
	Mode  string `json:"-"` // "create" or "update"
}

type ProgramTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (i *ProgramTypeInput) Validate() error {
	return shared.CustomValidate(i, map[string]string{
		"Name.required": "Program type name is required",
	})
}

func (i *ProgramTypeInput) ToProgramType() ProgramType {
	return ProgramType{
		ID:        i.ID,
		Name:      i.Name,
		CreatedBy: i.Actor,
		UpdatedBy: i.Actor,
	}
}

func (m *ProgramType) Update(newM *ProgramType) {
	m.Name = newM.Name
	m.UpdatedBy = newM.UpdatedBy
	m.UpdatedAt = time.Now()
}

func (m *ProgramType) UpdateFromInput(i ProgramTypeInput) {
	m.Name = i.Name
	m.UpdatedBy = i.Actor
	m.UpdatedAt = time.Now()
}

func (m *ProgramType) ToResponse() ProgramTypeResponse {
	return ProgramTypeResponse{
		ID:   m.ID,
		Name: m.Name,
	}
}

func ToProgramTypeResponses(models []ProgramType) []ProgramTypeResponse {
	res := make([]ProgramTypeResponse, len(models))
	for i, m := range models {
		res[i] = m.ToResponse()
	}
	return res
}
