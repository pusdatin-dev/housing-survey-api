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

type BalaiService interface {
	GetAll(ctx *fiber.Ctx) models.ServiceResponse
	GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse
	Create(ctx *fiber.Ctx, input *models.BalaiInput) models.ServiceResponse
	Update(ctx *fiber.Ctx, input *models.BalaiInput) models.ServiceResponse
	Delete(ctx *fiber.Ctx, id string) models.ServiceResponse
}

type balaiService struct {
	Db     *gorm.DB
	Config *config.Config
}

func NewBalaiService(ctx *context.AppContext) BalaiService {
	return &balaiService{
		Db:     ctx.DB,
		Config: ctx.Config,
	}
}

// ======= SERVICE METHODS =======

func (s *balaiService) GetAll(ctx *fiber.Ctx) models.ServiceResponse {
	var data []models.Balai
	db := s.Db.Model(&models.Balai{}).Where("deleted_at IS NULL")

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
		return models.InternalServerErrorResponse("Failed to count balais")
	}

	if err := db.Limit(limit).Offset(offset).Order("id ASC").Find(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve balais")
	}

	return models.OkResponse(http.StatusOK, "Success", fiber.Map{
		"data":       models.ToBalaiResponses(data),
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

func (s *balaiService) GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var data models.Balai
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Balai not found")
		}
		return models.InternalServerErrorResponse("Error retrieving balai")
	}
	return models.OkResponse(http.StatusOK, "Success", data.ToResponse())
}

func (s *balaiService) Create(ctx *fiber.Ctx, input *models.BalaiInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}
	data := input.ToModel()

	if err := s.Db.Create(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to create balai")
	}
	return models.OkResponse(http.StatusCreated, "Balai created", data.ToResponse())
}

func (s *balaiService) Update(ctx *fiber.Ctx, input *models.BalaiInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	var data models.Balai
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", input.ID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Balai not found")
		}
		return models.InternalServerErrorResponse("Error retrieving balai")
	}

	data.UpdateFromInput(input)

	if err := s.Db.Save(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to update balai")
	}
	return models.OkResponse(http.StatusOK, "Balai updated", data.ToResponse())
}

func (s *balaiService) Delete(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var data models.Balai
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse(fmt.Sprintf("Balai with id %s not found", id))
		}
		return models.InternalServerErrorResponse("Error retrieving balai")
	}

	data.MarkDeleted(utils.GetActor(ctx))
	if err := s.Db.Save(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to delete balai")
	}
	return models.OkResponse(http.StatusOK, "Balai deleted", nil)
}
