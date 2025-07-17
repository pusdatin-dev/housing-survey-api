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

type ProgramTypeService interface {
	GetAll(ctx *fiber.Ctx) models.ServiceResponse
	GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse
	Create(ctx *fiber.Ctx, input *models.ProgramTypeInput) models.ServiceResponse
	Update(ctx *fiber.Ctx, input *models.ProgramTypeInput) models.ServiceResponse
	Delete(ctx *fiber.Ctx, id string) models.ServiceResponse
}

type programTypeService struct {
	Db     *gorm.DB
	Config *config.Config
}

func NewProgramTypeService(ctx *context.AppContext) ProgramTypeService {
	return &programTypeService{
		Db:     ctx.DB,
		Config: ctx.Config,
	}
}

// ======= SERVICE METHODS =======

func (s *programTypeService) GetAll(ctx *fiber.Ctx) models.ServiceResponse {
	var data []models.ProgramType
	db := s.Db.Model(&models.ProgramType{}).Where("deleted_at IS NULL")

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
		return models.InternalServerErrorResponse("Failed to count program types")
	}

	if err := db.Limit(limit).Offset(offset).Order("id ASC").Find(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve program types")
	}

	return models.OkResponse(http.StatusOK, "Success", fiber.Map{
		"data":       models.ToProgramTypeResponses(data),
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

func (s *programTypeService) GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var data models.ProgramType
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("ProgramType not found")
		}
		return models.InternalServerErrorResponse("Error retrieving program type")
	}
	return models.OkResponse(http.StatusOK, "Success", data.ToResponse())
}

func (s *programTypeService) Create(ctx *fiber.Ctx, input *models.ProgramTypeInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}
	data := input.ToModel()

	if err := s.Db.Create(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to create program type")
	}
	return models.OkResponse(http.StatusCreated, "ProgramType created", data.ToResponse())
}

func (s *programTypeService) Update(ctx *fiber.Ctx, input *models.ProgramTypeInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	var data models.ProgramType
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", input.ID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("ProgramType not found")
		}
		return models.InternalServerErrorResponse("Error retrieving program type")
	}

	data.UpdateFromInput(input)

	if err := s.Db.Save(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to update program type")
	}
	return models.OkResponse(http.StatusOK, "ProgramType updated", data.ToResponse())
}

func (s *programTypeService) Delete(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var data models.ProgramType
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse(fmt.Sprintf("ProgramType with id %s not found", id))
		}
		return models.InternalServerErrorResponse("Error retrieving program type")
	}

	data.MarkDeleted(utils.GetActor(ctx))
	if err := s.Db.Save(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to delete program type")
	}
	return models.OkResponse(http.StatusOK, "ProgramType deleted", nil)
}
