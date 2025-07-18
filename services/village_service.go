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

type VillageService interface {
	GetAll(ctx *fiber.Ctx) models.ServiceResponse
	GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse
	Create(ctx *fiber.Ctx, input *models.VillageInput) models.ServiceResponse
	Update(ctx *fiber.Ctx, input *models.VillageInput) models.ServiceResponse
	Delete(ctx *fiber.Ctx, id string) models.ServiceResponse
}

type villageService struct {
	Db     *gorm.DB
	Config *config.Config
}

func NewVillageService(ctx *context.AppContext) VillageService {
	return &villageService{
		Db:     ctx.DB,
		Config: ctx.Config,
	}
}

// ======= SERVICE METHODS =======

func (s *villageService) GetAll(ctx *fiber.Ctx) models.ServiceResponse {
	var data []models.Village
	db := s.Db.Model(&models.Village{}).Where("deleted_at IS NULL")

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
		return models.InternalServerErrorResponse("Failed to count villages")
	}

	if err := db.Limit(limit).Offset(offset).Order("id ASC").Find(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve villages")
	}

	return models.OkResponse(http.StatusOK, "Success", fiber.Map{
		"data":       models.ToVillageResponses(data),
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

func (s *villageService) GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var data models.Village
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Village not found")
		}
		return models.InternalServerErrorResponse("Error retrieving village")
	}
	return models.OkResponse(http.StatusOK, "Success", data.ToResponse())
}

func (s *villageService) Create(ctx *fiber.Ctx, input *models.VillageInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}
	data := input.ToModel()

	if err := s.Db.FirstOrCreate(&data, &models.Village{ID: data.ID, Name: data.Name, SubdistrictID: data.SubdistrictID}).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to create village")
	}
	return models.OkResponse(http.StatusCreated, "Village created", data.ToResponse())
}

func (s *villageService) Update(ctx *fiber.Ctx, input *models.VillageInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	var data models.Village
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", input.ID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Village not found")
		}
		return models.InternalServerErrorResponse("Error retrieving village")
	}

	data.UpdateFromInput(input)

	if err := s.Db.Save(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to update village")
	}
	return models.OkResponse(http.StatusOK, "Village updated", data.ToResponse())
}

func (s *villageService) Delete(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var data models.Village
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse(fmt.Sprintf("Village with id %s not found", id))
		}
		return models.InternalServerErrorResponse("Error retrieving village")
	}

	data.MarkDeleted(utils.GetActor(ctx))
	if err := s.Db.Save(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to delete village")
	}
	return models.OkResponse(http.StatusOK, "Village deleted", nil)
}
