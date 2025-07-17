package models

import (
	"housing-survey-api/shared"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var _ = pq.StringArray{}

type Survey struct {
	ID                uint           `gorm:"primaryKey;autoIncrement"`
	UserID            uint           `gorm:"index"`
	Name              string         `gorm:"not null"`
	Address           string         `gorm:"not null"`
	Type              string         `gorm:"type:text;check:type IN ('Susun', 'Tapak');not null"`
	MbrStatus         string         `gorm:"type:text;check:mbr_status IN ('MBR', 'Non-MBR');not null"`
	Year              uint           `gorm:"index;not null"`
	UnitTarget        uint           `gorm:"index;not null"`
	StatusRealization string         `gorm:"check:status_realization IN ('Proses', 'Selesai')"`
	YearRealization   uint           `gorm:"index"`
	MonthRealization  uint           `gorm:"index"`
	Budget            uint64         // jumlah anggaran
	Coordinate        string         `gorm:"type:text"`
	StatusBalai       string         `gorm:"type:text;default:'Pending';check:status_balai IN ('Pending', 'Approved', 'Rejected')"` // Pending, Approved, Rejected
	StatusEselon1     string         `gorm:"type:text;default:'Pending';check:status_balai IN ('Pending', 'Approved', 'Rejected')"` // Pending, Approved, Rejected
	IsSubmitted       bool           `gorm:"default:false"`
	Notes             string         `gorm:"type:text"` // Notes for Balai or Eselon1
	ImagesBefore      pq.StringArray `gorm:"type:text[]"`
	ImagesAfter       pq.StringArray `gorm:"type:text[]"`
	ProvinceID        uint           `gorm:"index"`
	DistrictID        uint           `gorm:"index"`
	SubdistrictID     uint           `gorm:"index"`
	VillageID         uint           `gorm:"index"`
	User              User
	Province          Province
	District          District
	Subdistrict       Subdistrict
	Village           Village

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SurveyResponse struct {
	ID                uint           `json:"id"`
	UserID            uint           `json:"user_id"`
	UserEmail         string         `json:"user_email"`
	Name              string         `json:"survey_name"`
	Address           string         `json:"address"`
	Type              string         `json:"type"`
	MbrStatus         string         `json:"mbr_status"`
	Year              uint           `json:"year"`
	UnitTarget        uint           `json:"unit_target"`
	StatusRealization string         `json:"status_realization"`
	YearRealization   uint           `json:"year_realization"`
	MonthRealization  uint           `json:"month_realization"`
	Budget            uint64         `json:"budget"`
	Coordinate        string         `json:"coordinate"` // lat,lng string or GeoJSON
	Status            string         `json:"status"`
	StatusBalai       string         `json:"status_balai"`
	StatusEselon1     string         `json:"status_eselon1"`
	IsSubmitted       bool           `json:"is_submitted"` // default false
	Notes             string         `json:"notes"`
	ImagesBefore      pq.StringArray `json:"images_before"`
	ImagesAfter       pq.StringArray `json:"images_after"`
	ProvinceID        uint           `json:"province_id"`
	ProvinceName      string         `json:"province_name"`
	DistrictID        uint           `json:"district_id"`
	DistrictName      string         `json:"district_name"`
	SubdistrictID     uint           `json:"subdistrict_id"`
	SubdistrictName   string         `json:"subdistrict_name"`
	VillageID         uint           `json:"village_id"`
	VillageName       string         `json:"village_name"`
}

func (s *Survey) Update(newSurvey *Survey) {
	s.Name = newSurvey.Name
	s.Address = newSurvey.Address
	s.Type = newSurvey.Type
	s.MbrStatus = newSurvey.MbrStatus
	s.Year = newSurvey.Year
	s.UnitTarget = newSurvey.UnitTarget
	s.StatusRealization = newSurvey.StatusRealization
	s.YearRealization = newSurvey.YearRealization
	s.MonthRealization = newSurvey.MonthRealization
	s.Budget = newSurvey.Budget
	s.Coordinate = newSurvey.Coordinate
	s.IsSubmitted = newSurvey.IsSubmitted
	s.ImagesBefore = newSurvey.ImagesBefore
	s.ImagesAfter = newSurvey.ImagesAfter
	s.ProvinceID = newSurvey.ProvinceID
	s.DistrictID = newSurvey.DistrictID
	s.SubdistrictID = newSurvey.SubdistrictID
	s.VillageID = newSurvey.VillageID
	s.UpdatedBy = newSurvey.UpdatedBy
	s.UpdatedAt = time.Now()
}

func (s *Survey) UpdateFromInput(input SurveyInput) {
	var imagesBefore, imagesAfter pq.StringArray
	if s.StatusRealization == shared.StatusRealProses {
		imagesBefore = input.Images
	} else if s.StatusRealization == shared.StatusRealSelesai {
		imagesAfter = input.Images
	}
	s.ID = input.ID
	s.Name = input.Name
	s.Address = input.Address
	s.Type = input.Type
	s.MbrStatus = input.MbrStatus
	s.Year = input.Year
	s.UnitTarget = input.UnitTarget
	s.StatusRealization = input.StatusRealization
	s.YearRealization = input.YearRealization
	s.MonthRealization = input.MonthRealization
	s.Budget = input.Budget
	s.Coordinate = input.Coordinate
	s.IsSubmitted = input.IsSubmitted
	s.ImagesBefore = imagesBefore
	s.ImagesAfter = imagesAfter
	s.ProvinceID = input.ProvinceID
	s.DistrictID = input.DistrictID
	s.SubdistrictID = input.SubdistrictID
	s.VillageID = input.VillageID
	s.UpdatedBy = input.Actor
	s.UpdatedAt = time.Now()
}

func (s *Survey) ToResponse() SurveyResponse {
	return SurveyResponse{
		ID:                s.ID,
		UserID:            s.UserID,
		UserEmail:         s.User.Email,
		Name:              s.Name,
		Address:           s.Address,
		Type:              s.Type,
		MbrStatus:         s.MbrStatus,
		Year:              s.Year,
		UnitTarget:        s.UnitTarget,
		StatusRealization: s.StatusRealization,
		YearRealization:   s.YearRealization,
		MonthRealization:  s.MonthRealization,
		Budget:            s.Budget,
		Coordinate:        s.Coordinate,
		IsSubmitted:       s.IsSubmitted,
		Status:            s.GetStatusSurvey(),
		StatusBalai:       s.StatusBalai,
		StatusEselon1:     s.StatusEselon1,
		Notes:             s.Notes,
		ImagesBefore:      s.ImagesBefore,
		ImagesAfter:       s.ImagesAfter,
		ProvinceID:        s.ProvinceID,
		ProvinceName:      s.Province.Name,
		DistrictID:        s.DistrictID,
		DistrictName:      s.District.Name,
		SubdistrictID:     s.SubdistrictID,
		SubdistrictName:   s.Subdistrict.Name,
		VillageID:         s.VillageID,
		VillageName:       s.Village.Name,
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
	ID                uint           `json:"id"`
	UserID            uint           `json:"user_id" validate:"required"`
	Name              string         `json:"survey_name" validate:"required"`
	Address           string         `json:"address" validate:"required"`
	Type              string         `json:"type" validate:"required,oneof=Tapak Susun"`
	MbrStatus         string         `json:"mbr_status" validate:"required,oneof=MBR Non-MBR"`
	Year              uint           `json:"year" validate:"required"`
	UnitTarget        uint           `json:"unit_target" validate:"required"`
	StatusRealization string         `json:"status_realization" validate:"required,one of=Proses Selesai"`
	YearRealization   uint           `json:"year_realization"`
	MonthRealization  uint           `json:"month_realization"`
	Budget            uint64         `json:"budget"`
	Coordinate        string         `json:"coordinate"`   // lat,lng string or GeoJSON
	IsSubmitted       bool           `json:"is_submitted"` // default false
	Images            pq.StringArray `json:"images"`
	ProvinceID        uint           `json:"province_id" validate:"required"`
	DistrictID        uint           `json:"district_id" validate:"required"`
	SubdistrictID     uint           `json:"subdistrict_id" validate:"required"`
	VillageID         uint           `json:"village_id" validate:"required"`
	Actor             string         `json:"-"` // CreatedBy, UpdatedBy, DeletedBy
	Mode              string         `json:"-"` // "create" or "update"
}

// ToSurvey only used in creating survey
// if used for updating or viewing, adjust StatusBalai and StatusEselon1
func (s *SurveyInput) ToSurvey() Survey {
	var imagesBefore, imagesAfter pq.StringArray
	if s.StatusRealization == shared.StatusRealProses {
		imagesBefore = s.Images
	} else if s.StatusRealization == shared.StatusRealSelesai {
		imagesAfter = s.Images
	}
	var id uint
	if s.Mode == shared.Update {
		id = s.ID
	}
	return Survey{
		ID:               id,
		UserID:           s.UserID,
		Name:             s.Name,
		Type:             s.Type,
		MbrStatus:        s.MbrStatus,
		Year:             s.Year,
		UnitTarget:       s.UnitTarget,
		YearRealization:  s.YearRealization,
		MonthRealization: s.MonthRealization,
		Budget:           s.Budget,
		Coordinate:       s.Coordinate,
		IsSubmitted:      s.IsSubmitted,
		StatusBalai:      shared.Pending,
		StatusEselon1:    shared.Pending,
		ImagesBefore:     imagesBefore,
		ImagesAfter:      imagesAfter,
		ProvinceID:       s.ProvinceID,
		DistrictID:       s.DistrictID,
		SubdistrictID:    s.SubdistrictID,
		VillageID:        s.VillageID,
		CreatedBy:        s.Actor,
		CreatedAt:        time.Now(),
		UpdatedBy:        s.Actor,
		UpdatedAt:        time.Now(),
	}
}

func (s *SurveyInput) Validate() error {
	var customMessages = map[string]string{
		"ID.required":                "Survey ID is required for update",
		"UserID.required":            "User ID is required",
		"Name.required":              "Survey name is required",
		"Address.required":           "Address is required",
		"Type.required":              "Survey type is required",
		"Type.oneof":                 "Survey type must be either 'Tapak' or 'Susun'",
		"MBRStatus.required":         "MBR status is required",
		"MBRStatus.oneof":            "MBR status must be 'MBR' or 'Non-MBR'",
		"Year.required":              "Year is required",
		"UnitTarget.required":        "Unit target is required",
		"StatusRealization.required": "Status realization is required",
		"StatusRealization.oneof":    "Status realization must be either 'Proses' or 'Selesai'",
		"ProvinceID.required":        "Province is required",
		"DistrictID.required":        "District is required",
		"SubdistrictID.required":     "Subdistrict is required",
		"VillageID.required":         "Village is required",
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
