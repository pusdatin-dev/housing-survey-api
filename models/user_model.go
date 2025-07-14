package models

import (
	"housing-survey-api/shared"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email    string    `gorm:"uniqueIndex;not null"`
	Password string    `gorm:"not null"`
	Token    *string
	IsActive bool `gorm:"default:false"`
	RoleID   uint `gorm:"index"`
	Role     Role
	Profile  Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserResponse struct {
	ID       string    `json:"id"`
	Email    string    `json:"email"`
	IsActive bool      `json:"is_active"`
	RoleID   uint      `json:"role_id"`
	RoleName string    `json:"role_name"`
	Name     string    `json:"name"`
	BalaiID  uint      `json:"balai_id"`
	SKNo     string    `json:"sk_no"`
	SKDate   time.Time `json:"sk_date"`
	File     string    `json:"file"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:       u.ID.String(),
		Email:    u.Email,
		IsActive: u.IsActive,
		RoleID:   u.RoleID,
		RoleName: u.Role.Name,
		Name:     u.Profile.Name,
		BalaiID:  u.Profile.BalaiID,
		SKNo:     u.Profile.SKNo,
		SKDate:   u.Profile.SKDate,
		File:     u.Profile.File,
	}
}

func ToUserResponses(users []User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}
	return responses
}

type UserInput struct {
	ID       string    `json:"id"` // Optional, for updates. Should be a valid UUID.
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=6"`
	RoleID   uint      `json:"role_id" validate:"required"`
	Name     string    `json:"name"`
	BalaiID  uint      `json:"balai_id"`
	SKNo     string    `json:"sk_no"`
	SKDate   time.Time `json:"sk_date"`
	File     string    `json:"file"`
	Actor    string    `json:"-"`
}

// ToUser used only in creating a new user
// If wanted to use for updating, need to handle ID and Actor separately
func (u *UserInput) ToUser() User {
	return User{
		Email:    u.Email,
		Password: u.Password,
		RoleID:   u.RoleID,
		Profile: Profile{
			Name:      u.Name,
			BalaiID:   u.BalaiID,
			SKNo:      u.SKNo,
			SKDate:    u.SKDate,
			File:      u.File,
			CreatedBy: u.Actor,
			CreatedAt: time.Now(),
			UpdatedBy: u.Actor,
			UpdatedAt: time.Now(),
		},
		CreatedBy: u.Actor,
		CreatedAt: time.Now(),
		UpdatedBy: u.Actor,
		UpdatedAt: time.Now(),
	}
}

func (u *UserInput) Validate() error {
	customMessages := map[string]string{
		"Email.required":    "Email is required",
		"Email.email":       "Email must be a valid email address",
		"Password.required": "Password is required",
		"Password.min":      "Password must be at least 6 characters",
		"RoleID.required":   "Role is required",
	}
	return shared.CustomValidate(u, customMessages)
}

type ApprovingUser struct {
	UserIDs []string `json:"user_ids" validate:"required"`
	Actor   string   `json:"-"`
}

func (ap *ApprovingUser) Validate() error {
	customMessages := map[string]string{
		"UserIDs.required": "At least one user ID must be provided",
	}
	return shared.CustomValidate(ap, customMessages)
}
