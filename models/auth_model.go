package models

import (
	"housing-survey-api/shared"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (i *LoginInput) Validate() error {
	return shared.CustomValidate(i, map[string]string{
		"Email.required":    "Email is required",
		"Email.email":       "Invalid email format",
		"Password.required": "Password is required",
	})
}

type RefreshInput struct {
	RefreshToken string `json:"token" validate:"required"`
}

func (i *RefreshInput) Validate() error {
	return shared.CustomValidate(i, map[string]string{
		"RefreshToken.required": "Refresh token is required",
	})
}
