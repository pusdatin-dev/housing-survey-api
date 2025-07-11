package services

import (
	"housing-survey-api/config"
	"housing-survey-api/internal/context"
	"housing-survey-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CommentService interface {
	GetAllComments(ctx *fiber.Ctx) models.ServiceResponse
	CreatePublicComment(ctx *fiber.Ctx, input models.CommentInput) models.ServiceResponse
}

type commentService struct {
	Db     *gorm.DB
	Config *config.Config
}

func NewCommentService(ctx *context.AppContext) CommentService {
	return &commentService{
		Db:     ctx.DB,
		Config: ctx.Config,
	}
}

func (s *commentService) GetAllComments(ctx *fiber.Ctx) models.ServiceResponse {
	var comments []models.Comment
	db := s.Db.Model(&models.Comment{})

	// Filtering
	if name := ctx.Query("name"); name != "" {
		db = db.Where("name = ?", name)
	}
	if surveyId := ctx.Query("survey"); surveyId != "" {
		db = db.Where("survey_id = ?", surveyId)
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

	// Get paginated results
	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&comments).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve comments")
	}

	// Return with metadata
	return models.OkResponse(fiber.StatusOK, "Comments retrieved successfully", fiber.Map{
		"data":       comments,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int((total + int64(limit) - 1) / int64(limit)), // ceiling division
	})
}

func (s *commentService) CreatePublicComment(ctx *fiber.Ctx, input models.CommentInput) models.ServiceResponse {
	comment := input.ToComment()
	comment.CreatedBy = input.Actor

	if err := s.Db.Create(&comment).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to save comment")
	}

	return models.OkResponse(fiber.StatusOK, "Comment created", comment)
}
