package models

import (
	"housing-survey-api/shared"
	"time"

	"gorm.io/gorm"
)

type Resource struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	Name          string `gorm:"type:text;not null"`
	ProgramTypeID uint   `gorm:"index"`
	ProgramType   ProgramType
	CreatedBy     string `gorm:"type:text"`
	UpdatedBy     string `gorm:"type:text"`
	DeletedBy     string `gorm:"type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type ResourceInput struct {
	ID            uint   `json:"id"`
	Name          string `json:"name" validate:"required"`
	ProgramTypeID uint   `json:"program_type_id" validate:"required"`
	Actor         string `json:"-"`
	Mode          string `json:"-"`
}

type ResourceResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	ProgramTypeID uint   `json:"program_type_id"`
	ProgramType   string `json:"program_type"`
}

func (i *ResourceInput) Validate() error {
	return shared.CustomValidate(i, map[string]string{
		"Name.required":          "Resource name is required",
		"ProgramTypeID.required": "Program type ID is required",
	})
}

func (i *ResourceInput) ToResource() Resource {
	return Resource{
		ID:            i.ID,
		Name:          i.Name,
		ProgramTypeID: i.ProgramTypeID,
		CreatedBy:     i.Actor,
		UpdatedBy:     i.Actor,
	}
}

func (m *Resource) Update(newM *Resource) {
	m.Name = newM.Name
	m.ProgramTypeID = newM.ProgramTypeID
	m.UpdatedBy = newM.UpdatedBy
	m.UpdatedAt = time.Now()
}

func (m *Resource) UpdateFromInput(i ResourceInput) {
	m.Name = i.Name
	m.ProgramTypeID = i.ProgramTypeID
	m.UpdatedBy = i.Actor
	m.UpdatedAt = time.Now()
}

func (m *Resource) ToResponse() ResourceResponse {
	return ResourceResponse{
		ID:            m.ID,
		Name:          m.Name,
		ProgramTypeID: m.ProgramTypeID,
		ProgramType:   m.ProgramType.Name,
	}
}

func ToResourceResponses(models []Resource) []ResourceResponse {
	res := make([]ResourceResponse, len(models))
	for i, m := range models {
		res[i] = m.ToResponse()
	}
	return res
}
