package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var _ = pq.StringArray{}

type Comment struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	SurveyID uuid.UUID `gorm:"type:uuid;index"`
	Survey   Survey
	Name     string         `gorm:"not null"`
	Detail   string         `gorm:"type:text"`
	Images   pq.StringArray `gorm:"type:text[]"`

	CreatedBy string `gorm:"type:text"` // could be commenter name
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CommentResponse struct {
	ID        string         `json:"id"`
	SurveyID  string         `json:"survey_id"`
	Name      string         `json:"name"`
	Detail    string         `json:"detail"`
	Images    pq.StringArray `json:"images"`
	Address   string         `json:"address"`
	CreatedBy string         `json:"created_by"`
}

func (c *Comment) ToResponse() CommentResponse {
	return CommentResponse{
		ID:        c.ID.String(),
		SurveyID:  c.SurveyID.String(),
		Name:      c.Name,
		Detail:    c.Detail,
		Images:    c.Images,
		Address:   c.Survey.Address,
		CreatedBy: c.CreatedBy,
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
	SurveyID uuid.UUID `json:"survey_id"`
	Name     string    `json:"name"`
	Detail   string    `json:"detail"`
	Images   []string  `json:"images"`
	Actor    string    `json:"-"`
}

func (i CommentInput) ToComment() Comment {
	return Comment{
		SurveyID:  i.SurveyID,
		Name:      i.Name,
		Detail:    i.Detail,
		Images:    i.Images,
		CreatedBy: i.Actor,
		UpdatedBy: i.Actor,
	}
}
