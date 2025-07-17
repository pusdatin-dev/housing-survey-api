package models

import (
	"time"

	"housing-survey-api/shared"

	"gorm.io/gorm"
)

type Program struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"type:text;not null"`
	Detail     string
	ResourceID uint `gorm:"index"`
	Resource   Resource

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (p *Program) UpdateFromInput(input *ProgramInput) {
	p.Name = input.Name
	p.ResourceID = input.ResourceID
	p.UpdatedBy = input.Actor
	p.UpdatedAt = time.Now()
}

func (p *Program) Update(newProgram *Program) {
	p.Name = newProgram.Name
	p.ResourceID = newProgram.ResourceID
	p.UpdatedBy = newProgram.UpdatedBy
	p.UpdatedAt = time.Now()
}

func (p *Program) MarkDeleted(actor string) {
	p.DeletedBy = actor
	p.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
}

func (p *Program) ToResponse() ProgramResponse {
	return ProgramResponse{
		ID:         p.ID,
		Name:       p.Name,
		ResourceID: p.ResourceID,
		Resource:   p.Resource.Name,
	}
}

func ToProgramResponses(programs []Program) []ProgramResponse {
	responses := make([]ProgramResponse, len(programs))
	for i, program := range programs {
		responses[i] = program.ToResponse()
	}
	return responses
}

type ProgramResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	ResourceID uint   `json:"resource_id"`
	Resource   string `json:"resource_name"`
}

type ProgramInput struct {
	ID         uint   `json:"id"`
	Name       string `json:"name" validate:"required"`
	ResourceID uint   `json:"resource_id" validate:"required"`
	Actor      string `json:"-"` // created_by, updated_by
	Mode       string `json:"-"` // "create" or "update"
}

func (p *ProgramInput) Validate() error {
	custom := map[string]string{
		"Name.required":       "Program name is required",
		"ResourceID.required": "Resource is required",
	}
	return shared.CustomValidate(p, custom)
}

func (p *ProgramInput) ToModel() Program {
	now := time.Now()
	return Program{
		ID:         p.ID,
		Name:       p.Name,
		ResourceID: p.ResourceID,
		CreatedBy:  p.Actor,
		UpdatedBy:  p.Actor,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
