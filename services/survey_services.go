package services

import (
	"errors"
	"net/http"

	"housing-survey-api/config"
	"housing-survey-api/internal/context"
	"housing-survey-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SurveyService interface {
	GetAllSurveys(ctx *fiber.Ctx) models.ServiceResponse
	GetSurveyDetail(ctx *fiber.Ctx, id string) models.ServiceResponse
	CreateSurvey(ctx *fiber.Ctx, survey models.SurveyInput) models.ServiceResponse
	UpdateSurvey(ctx *fiber.Ctx, survey models.SurveyInput) models.ServiceResponse
	DeleteSurvey(ctx *fiber.Ctx, id string) models.ServiceResponse
}

type surveyService struct {
	Db     *gorm.DB
	Config *config.Config
}

func NewSurveyService(ctx *context.AppContext) SurveyService {
	return &surveyService{
		Db:     ctx.DB,
		Config: ctx.Config,
	}
}

func (s surveyService) GetAllSurveys(ctx *fiber.Ctx) (res models.ServiceResponse) {
	var surveys []models.Survey
	res.Status = true
	if err := s.Db.Order("created_at desc").Limit(100).Find(&surveys).Error; err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "Failed to retrieve surveys"
		return
	}
	res.Code = http.StatusOK
	res.Message = "Surveys retrieved successfully"
	res.Data = surveys
	return
}

func (s surveyService) GetSurveyDetail(ctx *fiber.Ctx, id string) (res models.ServiceResponse) {
	res.Status = true
	var survey models.Survey
	if err := s.Db.Where("id = ?", id).First(&survey).Error; err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "Failed to retrieve survey"
		return
	}
	if survey.ID == uuid.Nil {
		res.Code = http.StatusNotFound
		res.Message = "Survey not found"
		return
	}
	res.Code = http.StatusOK
	res.Message = "Survey retrieved successfully"
	res.Data = survey
	return
}

func (s surveyService) CreateSurvey(ctx *fiber.Ctx, survey models.SurveyInput) (res models.ServiceResponse) {
	// Extract role from Fiber context
	role, ok := ctx.Locals("role").(string)
	if !ok || role != "Surveyor" {
		res.Code = http.StatusForbidden
		res.Message = "Role not authorized to create survey"
		return
	}

	// Convert input to model
	newSurvey := survey.ToSurvey()
	newSurvey.ID = uuid.New() // Generate a new UUID for the survey

	// Insert into DB
	if err := s.Db.Create(&newSurvey).Error; err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "Failed to create survey"
		return
	}

	res.Status = true
	res.Code = http.StatusCreated
	res.Message = "Survey created successfully"
	res.Data = newSurvey
	return
}

func (s surveyService) UpdateSurvey(ctx *fiber.Ctx, survey models.SurveyInput) (res models.ServiceResponse) {
	// Extract role from Fiber context
	role, ok := ctx.Locals("role").(string)
	if !ok || role != "Surveyor" {
		res.Code = http.StatusForbidden
		res.Message = "Role not authorized to create survey"
		return
	}

	// Convert input to model
	newSurvey := survey.ToSurvey()
	oldSurvey := models.Survey{}
	if err := s.Db.Where("id = ?", survey.ID).First(&oldSurvey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.Code = http.StatusNotFound
			res.Message = "Survey not found"
			return
		}
		res.Code = http.StatusInternalServerError
		res.Message = "Failed to retrieve survey for update"
		return
	}

	// Insert into DB
	if err := s.Db.Create(&newSurvey).Error; err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "Failed to create survey"
		return
	}

	res.Status = true
	res.Code = http.StatusCreated
	res.Message = "Survey created successfully"
	res.Data = newSurvey
	return
}

func (s surveyService) DeleteSurvey(ctx *fiber.Ctx, id string) models.ServiceResponse {
	//TODO implement me
	return models.OkResponse(200, "Survey deleted successfully", nil)
}
