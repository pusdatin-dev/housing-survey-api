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

type ProvinceService interface {
	GetAll(ctx *fiber.Ctx) models.ServiceResponse
	GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse
	Create(ctx *fiber.Ctx, input *models.ProvinceInput) models.ServiceResponse
	Update(ctx *fiber.Ctx, input *models.ProvinceInput) models.ServiceResponse
	Delete(ctx *fiber.Ctx, id string) models.ServiceResponse
}

type provinceService struct {
	Db     *gorm.DB
	Config *config.Config
}

func NewProvinceService(ctx *context.AppContext) ProvinceService {
	return &provinceService{
		Db:     ctx.DB,
		Config: ctx.Config,
	}
}

// ======= SERVICE METHODS =======

func (s *provinceService) GetAll(ctx *fiber.Ctx) models.ServiceResponse {
	var provinces []models.Province
	db := s.Db.Model(&models.Province{}).Where("deleted_at IS NULL")

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
		return models.InternalServerErrorResponse("Failed to count provinces")
	}

	if err := db.Limit(limit).Offset(offset).Order("id ASC").Find(&provinces).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve provinces")
	}

	return models.OkResponse(http.StatusOK, "Success", fiber.Map{
		"data":       models.ToProvinceResponses(provinces),
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

func (s *provinceService) GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var province models.Province
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&province).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Province not found")
		}
		return models.InternalServerErrorResponse("Error retrieving province")
	}
	return models.OkResponse(http.StatusOK, "Success", province.ToResponse())
}

func (s *provinceService) Create(ctx *fiber.Ctx, input *models.ProvinceInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}
	province := input.ToModel()

	if err := s.Db.Create(&province).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to create province")
	}
	return models.OkResponse(http.StatusCreated, "Province created", province.ToResponse())
}

func (s *provinceService) Update(ctx *fiber.Ctx, input *models.ProvinceInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	var province models.Province
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", input.ID).First(&province).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Province not found")
		}
		return models.InternalServerErrorResponse("Error retrieving province")
	}

	province.UpdateFromInput(input)

	if err := s.Db.Save(&province).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to update province")
	}
	return models.OkResponse(http.StatusOK, "Province updated", province.ToResponse())
}

func (s *provinceService) Delete(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var province models.Province
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&province).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse(fmt.Sprintf("Province with id %s not found", id))
		}
		return models.InternalServerErrorResponse("Error retrieving province")
	}

	province.MarkDeleted(utils.GetActor(ctx))
	if err := s.Db.Save(&province).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to delete province")
	}
	return models.OkResponse(http.StatusOK, "Province deleted", nil)
}
