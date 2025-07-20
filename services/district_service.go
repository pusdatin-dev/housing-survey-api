package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"housing-survey-api/config"
	"housing-survey-api/internal/context"
	"housing-survey-api/models"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DistrictService interface {
	GetAll(ctx *fiber.Ctx) models.ServiceResponse
	GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse
	Create(ctx *fiber.Ctx, input *models.DistrictInput) models.ServiceResponse
	Update(ctx *fiber.Ctx, input *models.DistrictInput) models.ServiceResponse
	Delete(ctx *fiber.Ctx, id string) models.ServiceResponse
}

type districtService struct {
	Db     *gorm.DB
	Config *config.Config
}

func NewDistrictService(ctx *context.AppContext) DistrictService {
	return &districtService{
		Db:     ctx.DB,
		Config: ctx.Config,
	}
}

// ======= SERVICE METHODS =======

func (s *districtService) GetAll(ctx *fiber.Ctx) models.ServiceResponse {
	var districts []models.District
	db := s.Db.Model(&models.District{}).Where("deleted_at IS NULL")

	if search := ctx.Query("search"); search != "" {
		db = db.Where("name ILIKE ?", "%"+search+"%")
	}
	if province := ctx.Query("province"); province != "" {
		db = db.Where("province_id IN ?", strings.Split(province, ","))
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
		return models.InternalServerErrorResponse("Failed to count districts")
	}

	if err := db.Limit(limit).Offset(offset).Order("id ASC").Find(&districts).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve districts")
	}

	return models.OkResponse(http.StatusOK, "Success", fiber.Map{
		"data":       models.ToDistrictResponses(districts),
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

func (s *districtService) GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var district models.District
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&district).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("District not found")
		}
		return models.InternalServerErrorResponse("Error retrieving district")
	}
	return models.OkResponse(http.StatusOK, "Success", district.ToResponse())
}

func (s *districtService) Create(ctx *fiber.Ctx, input *models.DistrictInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}
	district := input.ToModel()

	if err := s.Db.FirstOrCreate(&district, &models.District{ID: district.ID, Name: district.Name, ProvinceID: district.ProvinceID}).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to create district")
	}
	return models.OkResponse(http.StatusCreated, "District created", district.ToResponse())
}

func (s *districtService) Update(ctx *fiber.Ctx, input *models.DistrictInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	var district models.District
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", input.ID).First(&district).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("District not found")
		}
		return models.InternalServerErrorResponse("Error retrieving district")
	}

	district.UpdateFromInput(input)

	if err := s.Db.Save(&district).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to update district")
	}
	return models.OkResponse(http.StatusOK, "District updated", district.ToResponse())
}

func (s *districtService) Delete(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var district models.District
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&district).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse(fmt.Sprintf("District with id %s not found", id))
		}
		return models.InternalServerErrorResponse("Error retrieving district")
	}

	district.MarkDeleted(utils.GetActor(ctx))
	if err := s.Db.Save(&district).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to delete district")
	}
	return models.OkResponse(http.StatusOK, "District deleted", nil)
}
