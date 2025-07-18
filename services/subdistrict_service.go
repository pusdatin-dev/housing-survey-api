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

type SubdistrictService interface {
	GetAll(ctx *fiber.Ctx) models.ServiceResponse
	GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse
	Create(ctx *fiber.Ctx, input *models.SubdistrictInput) models.ServiceResponse
	Update(ctx *fiber.Ctx, input *models.SubdistrictInput) models.ServiceResponse
	Delete(ctx *fiber.Ctx, id string) models.ServiceResponse
}

type subdistrictService struct {
	Db     *gorm.DB
	Config *config.Config
}

func NewSubdistrictService(ctx *context.AppContext) SubdistrictService {
	return &subdistrictService{
		Db:     ctx.DB,
		Config: ctx.Config,
	}
}

// ======= SERVICE METHODS =======

func (s *subdistrictService) GetAll(ctx *fiber.Ctx) models.ServiceResponse {
	var data []models.Subdistrict
	db := s.Db.Model(&models.Subdistrict{}).Where("deleted_at IS NULL")

	if search := ctx.Query("search"); search != "" {
		db = db.Where("name ILIKE ?", "%"+search+"%")
	}
	if district := ctx.Query("district"); district != "" {
		db = db.Where("district_id IN ?", strings.Split(district, ","))
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
		return models.InternalServerErrorResponse("Failed to count subdistricts")
	}

	if err := db.Limit(limit).Offset(offset).Order("id ASC").Find(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve subdistricts")
	}

	return models.OkResponse(http.StatusOK, "Success", fiber.Map{
		"data":       models.ToSubdistrictResponses(data),
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

func (s *subdistrictService) GetByID(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var data models.Subdistrict
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Subdistrict not found")
		}
		return models.InternalServerErrorResponse("Error retrieving subdistrict")
	}
	return models.OkResponse(http.StatusOK, "Success", data.ToResponse())
}

func (s *subdistrictService) Create(ctx *fiber.Ctx, input *models.SubdistrictInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}
	data := input.ToModel()

	if err := s.Db.FirstOrCreate(&data, &models.Subdistrict{ID: data.ID, Name: data.Name, DistrictID: data.DistrictID}).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to create subdistrict")
	}
	return models.OkResponse(http.StatusCreated, "Subdistrict created", data.ToResponse())
}

func (s *subdistrictService) Update(ctx *fiber.Ctx, input *models.SubdistrictInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	var data models.Subdistrict
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", input.ID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Subdistrict not found")
		}
		return models.InternalServerErrorResponse("Error retrieving subdistrict")
	}

	data.UpdateFromInput(input)

	if err := s.Db.Save(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to update subdistrict")
	}
	return models.OkResponse(http.StatusOK, "Subdistrict updated", data.ToResponse())
}

func (s *subdistrictService) Delete(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var data models.Subdistrict
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse(fmt.Sprintf("Subdistrict with id %s not found", id))
		}
		return models.InternalServerErrorResponse("Error retrieving subdistrict")
	}

	data.MarkDeleted(utils.GetActor(ctx))
	if err := s.Db.Save(&data).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to delete subdistrict")
	}
	return models.OkResponse(http.StatusOK, "Subdistrict deleted", nil)
}
