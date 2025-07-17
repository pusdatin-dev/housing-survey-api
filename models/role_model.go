package models

import (
	"time"

	"housing-survey-api/shared"

	"gorm.io/gorm"
)

type Role struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"uniqueIndex"`
	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type RoleInput struct {
	ID    uint   `json:"id"`
	Name  string `json:"name" validate:"required"`
	Actor string `json:"-"`
	Mode  string `json:"-"`
}

type RoleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (i *RoleInput) Validate() error {
	return shared.CustomValidate(i, map[string]string{
		"Name.required": "Role name is required",
	})
}

func (i *RoleInput) ToModel() Role {
	return Role{
		ID:        i.ID,
		Name:      i.Name,
		CreatedBy: i.Actor,
		UpdatedBy: i.Actor,
	}
}

func (m *Role) Update(newM *Role) {
	m.Name = newM.Name
	m.UpdatedBy = newM.UpdatedBy
	m.UpdatedAt = time.Now()
}

func (m *Role) UpdateFromInput(i *RoleInput) {
	m.Name = i.Name
	m.UpdatedBy = i.Actor
	m.UpdatedAt = time.Now()
}

func (m *Role) MarkDeleted(actor string) {
	m.DeletedBy = actor
	m.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
}

func (m *Role) ToResponse() RoleResponse {
	return RoleResponse{
		ID:   m.ID,
		Name: m.Name,
	}
}

func ToRoleResponses(models []Role) []RoleResponse {
	res := make([]RoleResponse, len(models))
	for i, m := range models {
		res[i] = m.ToResponse()
	}
	return res
}
