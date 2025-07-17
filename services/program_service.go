package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"housing-survey-api/config"
	"housing-survey-api/internal/context"
	"housing-survey-api/models"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProgramService interface {
	GetAll(ctx *fiber.Ctx) models.ServiceResponse
	GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse
	Create(ctx *fiber.Ctx, input *models.ProgramInput) models.ServiceResponse
	Update(ctx *fiber.Ctx, input *models.ProgramInput) models.ServiceResponse
	Delete(ctx *fiber.Ctx, id string) models.ServiceResponse
}

type programService struct {
	Db     *gorm.DB
	Config *config.Config
}

func NewProgramService(ctx *context.AppContext) ProgramService {
	return &programService{
		Db:     ctx.DB,
		Config: ctx.Config,
	}
}

// ======= SERVICE METHODS =======

func (s *programService) GetAll(ctx *fiber.Ctx) models.ServiceResponse {
	var data []models.Program
	db := s.Db.Model(&models.Program{}).Where("deleted_at IS NULL")

	if search := ctx.Query("search"); search != "" {
		db = db.Where("name ILIKE ?", "%"+search+"%")
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to count programs")
	}

	if err := db.Limit(limit).Offset(offset).Order("id ASC").Find(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve programs")
	}

	return models.OkResponse(http.StatusOK, "Success", fiber.Map{
		"data":       models.ToProgramResponses(data),
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

func (s *programService) GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var data models.Program
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Program not found")
		}
		return models.InternalServerErrorResponse("Error retrieving program")
	}
	return models.OkResponse(http.StatusOK, "Success", data.ToResponse())
}

func (s *programService) Create(ctx *fiber.Ctx, input *models.ProgramInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}
	data := input.ToModel()

	if err := s.Db.Create(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to create program")
	}
	return models.OkResponse(http.StatusCreated, "Program created", data.ToResponse())
}

func (s *programService) Update(ctx *fiber.Ctx, input *models.ProgramInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	var data models.Program
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", input.ID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Program not found")
		}
		return models.InternalServerErrorResponse("Error retrieving program")
	}

	data.UpdateFromInput(input)

	if err := s.Db.Save(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to update program")
	}
	return models.OkResponse(http.StatusOK, "Program updated", data.ToResponse())
}

func (s *programService) Delete(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var data models.Program
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse(fmt.Sprintf("Program with id %s not found", id))
		}
		return models.InternalServerErrorResponse("Error retrieving program")
	}

	data.MarkDeleted(utils.GetActor(ctx))
	if err := s.Db.Save(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to delete program")
	}
	return models.OkResponse(http.StatusOK, "Program deleted", nil)
}
