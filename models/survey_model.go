package models

import (
	"housing-survey-api/shared"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var _ = pq.StringArray{}

type Survey struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;index"`
	User          User
	Address       string         `gorm:"type:text"`
	Coordinate    string         `gorm:"type:text"` // lat,lng string or GeoJSON
	Type          string         `gorm:"type:text"`
	StatusBalai   string         `gorm:"type:text;default:'Pending';check:status_balai IN ('Pending', 'Approved', 'Rejected')"` // Pending, Approved, Rejected
	StatusEselon1 string         `gorm:"type:text;default:'Pending';check:status_balai IN ('Pending', 'Approved', 'Rejected')"` // Pending, Approved, Rejected
	IsSubmitted   bool           `gorm:"default:false"`
	Notes         string         `gorm:"type:text"` // Notes for Balai or Eselon1
	Images        pq.StringArray `gorm:"type:text[]"`
	ProvinceID    uint           `gorm:"index"`
	DistrictID    uint           `gorm:"index"`
	SubdistrictID uint           `gorm:"index"`
	VillageID     uint           `gorm:"index"`

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SurveyResponse struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	UserEmail     string         `json:"user_email"`
	Address       string         `json:"address"`
	Coordinate    string         `json:"coordinate"` // lat,lng string or GeoJSON
	Type          string         `json:"type"`
	IsSubmitted   bool           `json:"is_submitted"` // default false
	Status        string         `json:"status"`
	Notes         string         `json:"notes"`
	Images        pq.StringArray `json:"images"`
	ProvinceID    uint           `json:"province_id"`
	DistrictID    uint           `json:"district_id"`
	SubdistrictID uint           `json:"subdistrict_id"`
	VillageID     uint           `json:"village_id"`
}

func (s *Survey) Update(newSurvey *Survey) {
	s.Address = newSurvey.Address
	s.Coordinate = newSurvey.Coordinate
	s.Type = newSurvey.Type
	s.IsSubmitted = newSurvey.IsSubmitted
	s.Images = newSurvey.Images
	s.ProvinceID = newSurvey.ProvinceID
	s.DistrictID = newSurvey.DistrictID
	s.SubdistrictID = newSurvey.SubdistrictID
	s.VillageID = newSurvey.VillageID
	s.UpdatedBy = newSurvey.UpdatedBy
	s.UpdatedAt = time.Now()
}

func (s *Survey) UpdateFromInput(input SurveyInput) {
	s.Address = input.Address
	s.Coordinate = input.Coordinate
	s.Type = input.Type
	s.IsSubmitted = input.IsSubmitted
	s.Images = input.Images
	s.ProvinceID = input.ProvinceID
	s.DistrictID = input.DistrictID
	s.SubdistrictID = input.SubdistrictID
	s.VillageID = input.VillageID
	s.UpdatedBy = input.Actor
	s.UpdatedAt = time.Now()
}

func (s *Survey) ToResponse() SurveyResponse {
	return SurveyResponse{
		ID:            s.ID.String(),
		UserID:        s.UserID.String(),
		UserEmail:     s.User.Email,
		Address:       s.Address,
		Coordinate:    s.Coordinate,
		Type:          s.Type,
		IsSubmitted:   s.IsSubmitted,
		Status:        s.GetStatusSurvey(),
		Notes:         s.Notes,
		Images:        s.Images,
		ProvinceID:    s.ProvinceID,
		DistrictID:    s.DistrictID,
		SubdistrictID: s.SubdistrictID,
		VillageID:     s.VillageID,
	}
}

func (s *Survey) GetStatusSurvey() string {
	if !s.IsSubmitted {
		return shared.StatusDraft
	}
	if s.StatusBalai == shared.Pending && s.StatusEselon1 == shared.Pending {
		return shared.StatusWaitingBalai
	}
	if s.StatusBalai == shared.Approved && s.StatusEselon1 == shared.Pending {
		return shared.StatusWaitingEselon1
	}
	if s.StatusBalai == shared.Approved && s.StatusEselon1 == shared.Approved {
		return shared.StatusVerified
	}
	if s.StatusBalai == shared.Rejected {
		return shared.StatusRejectedBalai
	}
	if s.StatusEselon1 == shared.Rejected {
		return shared.StatusRejectedEselon1
	}
	return "unknown"
}

func ToSurveyResponse(surveys []Survey) []SurveyResponse {
	responses := make([]SurveyResponse, len(surveys))
	for i, survey := range surveys {
		responses[i] = survey.ToResponse()
	}
	return responses
}

type SurveyInput struct {
	ID            uuid.UUID      `json:"id"`
	UserID        uuid.UUID      `json:"user_id"`
	Address       string         `json:"address"`
	Coordinate    string         `json:"coordinate"` // lat,lng string or GeoJSON
	Type          string         `json:"type"`
	IsSubmitted   bool           `json:"is_submitted"` // default false
	Images        pq.StringArray `json:"images"`
	ProvinceID    uint           `json:"province_id"`
	DistrictID    uint           `json:"district_id"`
	SubdistrictID uint           `json:"subdistrict_id"`
	VillageID     uint           `json:"village_id"`
	Actor         string         `json:"-"` // CreatedBy, UpdatedBy, DeletedBy
	Mode          string         `json:"-"` // "create" or "update"
}

func (s *SurveyInput) ToSurvey() Survey {
	return Survey{
		ID:            s.ID,
		UserID:        s.UserID,
		Address:       s.Address,
		Coordinate:    s.Coordinate,
		Type:          s.Type,
		IsSubmitted:   s.IsSubmitted,
		StatusBalai:   shared.Pending,
		StatusEselon1: shared.Pending,
		Images:        s.Images,
		ProvinceID:    s.ProvinceID,
		DistrictID:    s.DistrictID,
		SubdistrictID: s.SubdistrictID,
		VillageID:     s.VillageID,
		CreatedBy:     s.Actor,
		CreatedAt:     time.Now(),
		UpdatedBy:     s.Actor,
		UpdatedAt:     time.Now(),
	}
}

func (s *SurveyInput) Validate() error {
	var customMessages = map[string]string{
		"ID.required":            "Survey ID is required for update",
		"UserID.required":        "User ID is required",
		"Address.required":       "Address is required",
		"Coordinate.required":    "Coordinate is required",
		"Type.required":          "Survey type is required",
		"ProvinceID.required":    "Province is required",
		"DistrictID.required":    "District is required",
		"SubdistrictID.required": "Subdistrict is required",
		"VillageID.required":     "Village is required",
	}

	return shared.CustomValidate(s, customMessages)
}

type SurveyActionInput struct {
	SurveyIDs []string `json:"survey_ids" validate:"required"`
	Action    string   `json:"action" validate:"required,oneof=Approved Rejected"`
	Notes     string   `json:"notes" validation:"required_if=Action Rejected"` // Notes for rejection
	Actor     string   `json:"-"`                                              // Actor who performs the action
}

func (s *SurveyActionInput) Validate() error {
	customMessages := map[string]string{
		"SurveyIDs.required": "Survey IDs are required",
		"Action.required":    "Action is required",
		"Action.oneof":       "Action must be either 'Approved' or 'Rejected'",
		"Notes.required_if":  "Notes are required when rejecting a survey",
	}
	return shared.CustomValidate(s, customMessages)
}
