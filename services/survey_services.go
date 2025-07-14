package services

import (
	"errors"
	"fmt"
	"housing-survey-api/config"
	"housing-survey-api/internal/context"
	"housing-survey-api/models"
	"housing-survey-api/shared"
	"housing-survey-api/utils"
	"strconv"
	"time"

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
	ActionSurvey(ctx *fiber.Ctx, input models.SurveyActionInput) models.ServiceResponse
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

func (s *surveyService) GetAllSurveys(ctx *fiber.Ctx) models.ServiceResponse {
	var surveys []models.Survey
	db := s.Db.Model(&models.Survey{}).Where("deleted_at IS NULL")

	// Filtering
	if address := ctx.Query("address"); address != "" {
		db = db.Where("address LIKE ?", "%"+address+"%")
	}
	if userId := ctx.Query("user_id"); userId != "" {
		if _, err := uuid.Parse(userId); err != nil {
			return models.BadRequestResponse("Invalid user ID format")
		}
		db = db.Where("user_id = ?", userId)
	}
	if types := ctx.Query("types"); types != "" {
		// Assuming types is a comma-separated list of survey types
		typeList := utils.SplitAndTrim(types, ",")
		if len(typeList) > 0 {
			db = db.Where("type IN ?", typeList)
		}
	}
	if provinceIDs := ctx.Query("province_ids"); provinceIDs != "" {
		// Assuming province_ids is a comma-separated list of province IDs
		provinceIDList := utils.SplitAndTrim(provinceIDs, ",")
		if len(provinceIDList) > 0 {
			db = db.Where("province_id IN ?", provinceIDList)
		}
	}
	if districtIDs := ctx.Query("district_ids"); districtIDs != "" {
		// Assuming district_ids is a comma-separated list of district IDs
		districtIDList := utils.SplitAndTrim(districtIDs, ",")
		if len(districtIDList) > 0 {
			db = db.Where("district_id IN ?", districtIDList)
		}
	}
	if subdistrictIDs := ctx.Query("subdistrict_ids"); subdistrictIDs != "" {
		// Assuming subdistrict_ids is a comma-separated list of subdistrict IDs
		subdistrictIDList := utils.SplitAndTrim(subdistrictIDs, ",")
		if len(subdistrictIDList) > 0 {
			db = db.Where("subdistrict_id IN ?", subdistrictIDList)
		}
	}
	if villageIDs := ctx.Query("village_ids"); villageIDs != "" {
		// Assuming village_ids is a comma-separated list of village IDs
		villageIDList := utils.SplitAndTrim(villageIDs, ",")
		if len(villageIDList) > 0 {
			db = db.Where("village_id IN ?", villageIDList)
		}
	}

	// Pagination
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Count total results
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to count comments")
	}

	if err := s.Db.Preload("User").Limit(limit).Offset(offset).Order("created_at desc").Find(&surveys).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve surveys")
	}
	// Return with metadata
	return models.OkResponse(fiber.StatusOK, "Survey retrieved successfully", fiber.Map{
		"data":       models.ToSurveyResponse(surveys),
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int((total + int64(limit) - 1) / int64(limit)), // ceiling division
	})
}

func (s *surveyService) GetSurveyDetail(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var survey models.Survey
	if err := s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&survey).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve survey")
	}
	if survey.ID == uuid.Nil {
		return models.NotFoundResponse("Survey not found")
	}
	return models.OkResponse(fiber.StatusOK, "Survey retrieved successfully", survey.ToResponse())
}

func (s *surveyService) CreateSurvey(ctx *fiber.Ctx, input models.SurveyInput) models.ServiceResponse {
	//enforcing role surveyor only will be in middleware
	// Convert input to model
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}
	survey := input.ToSurvey()
	survey.ID = uuid.New() // Generate a new UUID for the survey

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Cannot find UserID in token")
	}

	if userID != survey.UserID.String() {
		return models.BadRequestResponse("Cannot create survey for another user")
	}

	// Insert into DB
	if err := s.Db.Create(&survey).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to create survey")
	}

	return models.OkResponse(fiber.StatusCreated, "Survey created successfully", survey.ToResponse())
}

func (s *surveyService) UpdateSurvey(ctx *fiber.Ctx, survey models.SurveyInput) models.ServiceResponse {
	//enforcing role surveyor only will be in middleware
	//newSurvey := survey.ToSurvey()
	oldSurvey := models.Survey{}

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Cannot find UserID in token")
	}

	if userID != survey.UserID.String() {
		return models.BadRequestResponse("Cannot update survey for another user")
	}

	if err := s.Db.Where("id = ? AND deleted_at IS NULL", survey.ID).First(&oldSurvey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Survey not found")
		}
		return models.InternalServerErrorResponse("Failed to retrieve survey for update")
	}

	// Insert into DB
	oldSurvey.UpdateFromInput(survey)
	if err := s.Db.Save(&oldSurvey).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to update survey")
	}

	return models.OkResponse(fiber.StatusCreated, "Survey created successfully", oldSurvey.ToResponse())
}

func (s *surveyService) DeleteSurvey(ctx *fiber.Ctx, id string) models.ServiceResponse {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Cannot find UserID in token")
	}

	var survey models.Survey
	if err = s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&survey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse(fmt.Sprintf("Survey with id %s not found", id))
		}
		return models.InternalServerErrorResponse(fmt.Sprintf("Failed to retrieve survey with id %s", id))
	}

	if userID != survey.UserID.String() {
		return models.BadRequestResponse("Cannot delete survey for another user")
	}

	survey.DeletedBy = userID
	survey.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
	if err = s.Db.Save(&survey).Error; err != nil {
		return models.InternalServerErrorResponse(fmt.Sprintf("Failed to delete survey with id %s", id))
	}

	return models.OkResponse(200, "Survey deleted successfully", nil)
}

func (s *surveyService) ActionSurvey(ctx *fiber.Ctx, input models.SurveyActionInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	role, err := utils.GetRoleNameFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Cannot determine role")
	}

	isVerificatorBalai := role == s.Config.Roles.VerificatorBalai
	isVerificatorEselon1 := role == s.Config.Roles.VerificatorEselon1

	if !isVerificatorBalai && !isVerificatorEselon1 {
		return models.ForbiddenResponse("You are not allowed to perform this action")
	}

	// Base query
	db := s.Db.Model(&models.Survey{}).
		Where("id IN ?", input.SurveyIDs).
		Where("is_submitted = ? AND deleted_at IS NULL", true)

	// Filter based on role
	if isVerificatorBalai {
		db = db.Where("status_balai = ?", shared.Pending)
	} else if isVerificatorEselon1 {
		db = db.Where("status_balai = ? AND status_eselon1 = ?", shared.Approved, shared.Pending)
	}

	// Prepare update map
	update := map[string]interface{}{}
	if input.Action == shared.Rejected {
		update["notes"] = input.Notes
	}
	if isVerificatorBalai {
		update["status_balai"] = input.Action
	} else {
		update["status_eselon1"] = input.Action
	}

	// Perform update in one query
	result := db.Updates(update)
	if result.Error != nil {
		return models.InternalServerErrorResponse("Failed to update surveys")
	}

	// Calculate counts
	successCount := result.RowsAffected
	failedCount := int64(len(input.SurveyIDs)) - successCount

	return models.OkResponse(fiber.StatusOK, fmt.Sprintf(
		"%s %d survey(s), %d failed", input.Action, successCount, failedCount,
	), fiber.Map{
		"success_count": successCount,
		"failed_count":  failedCount,
	})
}
