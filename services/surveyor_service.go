package services

import (
	"housing-survey-api/internal/context"
	"housing-survey-api/models"
)

type SurveyorService interface {
	GetAllSurveyors() ([]models.Surveyor, error)
}

type surveyorService struct {
	ctx *context.AppContext // simpan context, akses ke DB/Config dll
}

// Constructor NewSurveyorService: dependency injection pakai AppContext
func NewSurveyorService(ctx *context.AppContext) SurveyorService {
	return &surveyorService{ctx: ctx}
}

// Implementasi fungsi GetAllSurveyors
func (s *surveyorService) GetAllSurveyors() ([]models.Surveyor, error) {
	var surveyors []models.Surveyor

	// Query database lewat s.ctx.DB
	if err := s.ctx.DB.Find(&surveyors).Error; err != nil {
		return nil, err
	}
	return surveyors, nil
}
