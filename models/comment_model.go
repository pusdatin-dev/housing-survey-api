package models

import (
	"housing-survey-api/shared"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var _ = pq.StringArray{}

type Comment struct {
	ID         uint           `gorm:"primaryKey;autoIncrement"`
	UserID     *uint          `gorm:"index"` // nullable, if it's a public comment
	SurveyID   uint           `gorm:"index"`
	ParentID   uint           `gorm:"index"`
	IsResolved bool           `gorm:"default:false"`
	Name       string         `gorm:"not null"`
	Detail     string         `gorm:"type:text"`
	Images     pq.StringArray `gorm:"type:text[]"`
	Survey     Survey         `gorm:"foreignKey:SurveyID"`
	Children   []Comment      `gorm:"-"`

	ResolvedBy *string
	ResolvedAt *time.Time
	CreatedBy  string `gorm:"type:text"` // could be commenter name
	UpdatedBy  string `gorm:"type:text"`
	DeletedBy  string `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (c *Comment) UpdateFromInput(input CommentInput) {
	c.Name = input.Name
	c.Detail = input.Detail
	c.IsResolved = input.IsResolved
	c.Images = pq.StringArray(input.Images)
	c.UpdatedBy = input.Actor
	c.UpdatedAt = time.Now()
}

func (c *Comment) MarkDeleted(actor string) {
	c.DeletedBy = actor
	c.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
}

func (c *Comment) MarkAction(action, actor string) {
	c.IsResolved = false
	c.ResolvedBy = nil
	c.ResolvedAt = nil
	if action == shared.StatusResolved {
		c.IsResolved = true
		c.ResolvedBy = &actor
		resolvedAt := time.Now()
		c.ResolvedAt = &resolvedAt
	}
}

type CommentResponse struct {
	ID         uint              `json:"id"`
	UserID     uint              `json:"user_id"` // nullable, if it's a public comment
	SurveyID   uint              `json:"survey_id"`
	ParentID   uint              `json:"parent_id"`
	SurveyName string            `json:"survey_name"`
	Name       string            `json:"name"`
	Detail     string            `json:"detail"`
	IsResolved bool              `json:"is_resolved"`
	Children   []CommentResponse `json:"children,omitempty"`
	Images     pq.StringArray    `json:"images"`
	CreatedBy  string            `json:"created_by"`
	CreatedAdt time.Time         `json:"created_at"`
	ResolvedBy string            `json:"resolved_by"`
	ResolvedAt *time.Time        `json:"resolved_at"`
}

func (c *Comment) ToResponse() CommentResponse {
	userID := uint(0)
	resolvedBy := ""
	if c.UserID != nil {
		userID = *c.UserID
	}
	if c.ResolvedBy != nil {
		resolvedBy = *c.ResolvedBy
	}

	children := make([]CommentResponse, len(c.Children))
	for i := range c.Children {
		children[i] = c.Children[i].ToResponse()
	}

	return CommentResponse{
		ID:         c.ID,
		UserID:     userID,
		SurveyID:   c.SurveyID,
		ParentID:   c.ParentID,
		SurveyName: c.Survey.Name,
		Name:       c.Name,
		Detail:     c.Detail,
		IsResolved: c.IsResolved,
		Images:     c.Images,
		CreatedBy:  c.CreatedBy,
		CreatedAdt: c.CreatedAt,
		ResolvedBy: resolvedBy,
		ResolvedAt: c.ResolvedAt,
		Children:   children,
	}
}

func ToCommentResponses(comments []Comment) []CommentResponse {
	responses := make([]CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = comment.ToResponse()
	}
	return responses
}

type CommentInput struct {
	ID         uint     `json:"id" validate:"required_if=Mode update"`
	UserID     *uint    `json:"user_id"` // nullable, if it's a public comment
	SurveyID   uint     `json:"survey_id" validate:"required"`
	ParentID   uint     `json:"parent_id"`
	Name       string   `json:"name" validate:"required"`
	Detail     string   `json:"detail" validate:"required"`
	IsResolved bool     `json:"is_resolved"`
	Images     []string `json:"images"`
	Actor      string   `json:"-"`
	Mode       string   `json:"-"`
}

func (i CommentInput) ToComment() Comment {
	return Comment{
		ID:         i.ID,
		SurveyID:   i.SurveyID,
		ParentID:   i.ParentID,
		Name:       i.Name,
		Detail:     i.Detail,
		IsResolved: i.IsResolved,
		Images:     i.Images,
		CreatedBy:  i.Actor,
		CreatedAt:  time.Now(),
		UpdatedBy:  i.Actor,
		UpdatedAt:  time.Now(),
	}
}

func (i CommentInput) Validate() error {
	return shared.CustomValidate(i, map[string]string{
		"ID.required_if":    "ID is required when updating a comment",
		"SurveyID.required": "Survey ID is required",
		"Name.required":     "Name is required",
		"Detail.required":   "Detail is required",
	})
}

type CommentActionInput struct {
	CommentID uint   `json:"comment_id" validate:"required"`
	Action    string `json:"action" validate:"required,oneof=Resolved Unresolved"`
	Actor     string `json:"-"` // Actor who performs the action
}

func (i CommentActionInput) Validate() error {
	return shared.CustomValidate(i, map[string]string{
		"CommentID.required": "Comment ID is required",
		"Action.required":    "Action is required",
		"Action.oneof":       "Action must be either 'Resolved' or 'Unresolved'",
	})
}
