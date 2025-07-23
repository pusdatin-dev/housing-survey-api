package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"housing-survey-api/config"
	"housing-survey-api/internal/context"
	"housing-survey-api/models"
	"housing-survey-api/shared"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CommentService interface {
	GetAllComments(ctx *fiber.Ctx) models.ServiceResponse
	GetCommentByID(ctx *fiber.Ctx) models.ServiceResponse
	CreatePublicComment(ctx *fiber.Ctx, input models.CommentInput) models.ServiceResponse
	UpdateComment(ctx *fiber.Ctx, input models.CommentInput) models.ServiceResponse
	DeleteComment(ctx *fiber.Ctx, id string) models.ServiceResponse
	ActionComment(ctx *fiber.Ctx, input models.CommentActionInput) models.ServiceResponse
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
	var allComments []models.Comment

	db := s.Db.Preload("Survey").Where("deleted_at IS NULL")

	// Filtering
	if surveyId := ctx.Query("survey"); surveyId != "" {
		db = db.Where("survey_id = ?", surveyId)
	}
	if keyword := ctx.Query("q"); keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("name ILIKE ? OR detail ILIKE ?", like, like)
	}

	// Fetch all comments (we need all to build the full tree)
	if err := db.Order("created_at ASC").Find(&allComments).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve comments")
	}

	// Build full comment tree
	rootComments := BuildCommentTree(allComments)

	// Pagination
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	total := len(rootComments)

	// Apply pagination to root-level comments only
	end := offset + limit
	if end > total {
		end = total
	}
	paginatedRoots := rootComments
	if offset < total {
		paginatedRoots = rootComments[offset:end]
	} else {
		paginatedRoots = []models.Comment{}
	}

	// Convert to response
	response := make([]models.CommentResponse, len(paginatedRoots))
	for i := range paginatedRoots {
		response[i] = paginatedRoots[i].ToResponse()
	}

	return models.OkResponse(fiber.StatusOK, "Comments retrieved successfully", fiber.Map{
		"data":       response,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int((int64(total) + int64(limit) - 1) / int64(limit)),
	})
}

func BuildCommentTree(comments []models.Comment) []models.Comment {
	// Map parent_id -> []comments
	childrenMap := make(map[uint][]models.Comment)
	var rootComments []models.Comment

	// Group all comments by ParentID
	for _, c := range comments {
		childrenMap[c.ParentID] = append(childrenMap[c.ParentID], c)
	}

	// Recursively attach children
	var attachChildren func(comment *models.Comment)
	attachChildren = func(comment *models.Comment) {
		comment.Children = childrenMap[comment.ID]
		for i := range comment.Children {
			attachChildren(&comment.Children[i])
		}
	}

	// Start from comments with ParentID == 0
	for _, root := range childrenMap[0] {
		attachChildren(&root)
		rootComments = append(rootComments, root)
	}

	return rootComments
}

func (s *commentService) GetCommentByID(ctx *fiber.Ctx) models.ServiceResponse {
	id := ctx.Params("id")
	if id == "" {
		return models.BadRequestResponse("Comment ID is required")
	}

	var comment models.Comment
	if err := s.Db.Preload("Survey").Where("id = ?", id).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Comment not found")
		}
		return models.InternalServerErrorResponse("Failed to retrieve comment")
	}

	return models.OkResponse(fiber.StatusOK, "Comment retrieved successfully", comment.ToResponse())
}

func (s *commentService) CreatePublicComment(ctx *fiber.Ctx, input models.CommentInput) models.ServiceResponse {
	if s.ContainsInappropriateContent(input.Name) || s.ContainsInappropriateContent(input.Detail) {
		return models.BadRequestResponse("Inappropriate content detected.")
	}
	comment := input.ToComment()

	var survey models.Survey
	if err := s.Db.Where("id = ?", input.SurveyID).First(&survey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Survey not found")
		}
		return models.InternalServerErrorResponse("Failed to find survey")
	}

	for _, image := range comment.Images {
		if !utils.IsValidBase64Image(image) {
			return models.BadRequestResponse("Invalid image format in base64")
		}
	}

	if err := s.Db.Create(&comment).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to save comment")
	}
	comment.Survey = survey
	return models.OkResponse(fiber.StatusOK, "Comment created", comment.ToResponse())
}

func (s *commentService) ContainsInappropriateContent(input string) bool {
	lower := strings.ToLower(input)
	for _, word := range s.Config.BannedWords {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}

func (s *commentService) UpdateComment(ctx *fiber.Ctx, input models.CommentInput) models.ServiceResponse {
	action := "UPDATE_COMMENT"
	var userID int
	var err error
	if err = input.Validate(); err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.BadRequestResponse(err.Error())
	}

	if userID, err = utils.GetUserIDFromContext(ctx); err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.UnauthorizedResponse("User not authenticated")
	}

	if uint(userID) != *input.UserID {
		utils.LogAudit(ctx, action, "Cannot update another user's comment")
		return models.ForbiddenResponse("Cannot update another user's comment")
	}
	if s.ContainsInappropriateContent(input.Name) || s.ContainsInappropriateContent(input.Detail) {
		utils.LogAudit(ctx, action, "Inappropriate content detected in comment")
		return models.BadRequestResponse("Inappropriate content detected.")
	}

	var comment models.Comment
	if err = s.Db.Where("id = ?", input.ID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogAudit(ctx, action, err.Error())
			return models.NotFoundResponse("Comment not found")
		}
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Failed to retrieve comment")
	}

	comment.UpdateFromInput(input)
	if err = s.Db.Save(&comment).Error; err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Failed to update comment")
	}

	return models.OkResponse(fiber.StatusOK, "Comment updated", comment.ToResponse())
}

func (s *commentService) DeleteComment(ctx *fiber.Ctx, id string) models.ServiceResponse {
	action := "DELETE_COMMENT"
	if id == "" {
		utils.LogAudit(ctx, action, "Comment ID is required")
		return models.BadRequestResponse("Comment ID is required")
	}

	actor := utils.GetActor(ctx)
	var userID int
	var err error
	if userID, err = utils.GetUserIDFromContext(ctx); err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.UnauthorizedResponse("User not authenticated")
	}

	var comment models.Comment
	if err := s.Db.Where("id = ?", id).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogAudit(ctx, action, err.Error())
			return models.NotFoundResponse("Comment not found")
		}
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Failed to retrieve comment")
	}

	if comment.UserID != nil && *comment.UserID != uint(userID) {
		utils.LogAudit(ctx, action, "User does not have permission to delete this comment")
		return models.ForbiddenResponse("You do not have permission to delete this comment")
	}

	comment.MarkDeleted(actor)

	if err := s.Db.Save(&comment).Error; err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Failed to delete comment")
	}

	return models.OkResponse(fiber.StatusOK, "Comment deleted successfully", nil)
}

func (s *commentService) ActionComment(ctx *fiber.Ctx, input models.CommentActionInput) models.ServiceResponse {
	actionTag := "ACTION_COMMENT"
	if err := input.Validate(); err != nil {
		utils.LogAudit(ctx, actionTag, err.Error())
		return models.BadRequestResponse(err.Error())
	}

	actorID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		utils.LogAudit(ctx, actionTag, err.Error())
		return models.UnauthorizedResponse("User not authenticated")
	}
	actorIDStr := fmt.Sprint(actorID)

	actorRole, err := utils.GetRoleNameFromContext(ctx)
	if err != nil {
		utils.LogAudit(ctx, actionTag, err.Error())
		return models.UnauthorizedResponse("User role not found")
	}

	var root models.Comment
	if err := s.Db.Where("id = ?", input.CommentID).First(&root).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogAudit(ctx, actionTag, err.Error())
			return models.NotFoundResponse("Comment not found")
		}
		utils.LogAudit(ctx, actionTag, err.Error())
		return models.InternalServerErrorResponse("Failed to retrieve comment")
	}

	// Check permission: either comment owner or PIC role
	fmt.Println("actorRole:", actorRole, "root.UserID:", root.UserID, "actorID:", actorID, root.UserID == nil || *root.UserID != uint(actorID) && !shared.PICSurvey[actorRole])
	if (root.UserID == nil || *root.UserID != uint(actorID)) && !shared.PICSurvey[actorRole] {
		utils.LogAudit(ctx, actionTag, "Forbidden actionTag")
		return models.ForbiddenResponse("You do not have permission to perform this actionTag")
	}

	// Run update in a transaction
	if err := s.Db.Transaction(func(tx *gorm.DB) error {
		var allRelated []models.Comment

		// Get all descendants using recursive CTE
		if err := tx.Raw(`
			WITH RECURSIVE descendants AS (
				SELECT * FROM comments WHERE id = ?
				UNION ALL
				SELECT c.* FROM comments c
				INNER JOIN descendants d ON c.parent_id = d.id
			)
			SELECT * FROM descendants;
		`, root.ID).Scan(&allRelated).Error; err != nil {
			return err
		}

		// Get all ancestors using recursive CTE (if applicable)
		if root.ParentID != 0 {
			var ancestors []models.Comment
			if err := tx.Raw(`
				WITH RECURSIVE ancestors AS (
					SELECT * FROM comments WHERE id = ?
					UNION ALL
					SELECT c.* FROM comments c
					INNER JOIN ancestors a ON a.parent_id = c.id
				)
				SELECT * FROM ancestors;
			`, root.ParentID).Scan(&ancestors).Error; err != nil {
				return err
			}
			allRelated = append(allRelated, ancestors...)
		}

		// Mark all affected comments
		for i := range allRelated {
			allRelated[i].MarkAction(input.Action, actorIDStr)
			if err := tx.Save(&allRelated[i]).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		utils.LogAudit(ctx, actionTag, err.Error())
		return models.InternalServerErrorResponse("Failed to update comment status")
	}

	return models.OkResponse(fiber.StatusOK, "Action performed successfully", nil)
}
