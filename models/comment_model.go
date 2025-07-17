package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var _ = pq.StringArray{}

type Comment struct {
	ID         uint           `gorm:"primaryKey;autoIncrement"`
	SurveyID   uint           `gorm:"index"`
	ParentID   uint           `gorm:"index"`
	IsResolved bool           `gorm:"default:false"`
	Name       string         `gorm:"not null"`
	Detail     string         `gorm:"type:text"`
	Images     pq.StringArray `gorm:"type:text[]"`
	Survey     Survey

	CreatedBy string `gorm:"type:text"` // could be commenter name
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CommentResponse struct {
	ID         uint           `json:"id"`
	SurveyID   uint           `json:"survey_id"`
	ParentID   uint           `json:"parent_id"`
	SurveyName string         `json:"survey_name"`
	Name       string         `json:"name"`
	Detail     string         `json:"detail"`
	IsResolved bool           `json:"is_resolved"`
	Images     pq.StringArray `json:"images"`
	CreatedBy  string         `json:"created_by"`
}

func (c *Comment) ToResponse() CommentResponse {
	return CommentResponse{
		ID:         c.ID,
		SurveyID:   c.SurveyID,
		ParentID:   c.ParentID,
		SurveyName: c.Survey.Name,
		Name:       c.Name,
		Detail:     c.Detail,
		IsResolved: c.IsResolved,
		Images:     c.Images,
		CreatedBy:  c.CreatedBy,
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
	SurveyID   uint     `json:"survey_id"`
	ParentID   uint     `json:"parent_id"`
	Name       string   `json:"name"`
	Detail     string   `json:"detail"`
	IsResolved bool     `json:"is_resolved"`
	Images     []string `json:"images"`
	Actor      string   `json:"-"`
}

func (i CommentInput) ToComment() Comment {
	return Comment{
		SurveyID:   i.SurveyID,
		ParentID:   i.ParentID,
		Name:       i.Name,
		Detail:     i.Detail,
		IsResolved: i.IsResolved,
		Images:     i.Images,
		CreatedBy:  i.Actor,
		UpdatedBy:  i.Actor,
	}
}
