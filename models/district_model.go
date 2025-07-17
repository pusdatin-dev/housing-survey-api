package models

import (
	"time"

	"housing-survey-api/shared"

	"gorm.io/gorm"
)

//
// ====== Model ======
//

// District Master Data (Kabupaten/Kota)
type District struct {
	ID         uint     `gorm:"primaryKey;autoIncrement"`
	Name       string   `gorm:"type:text;uniqueIndex;not null"`
	ProvinceID uint     `gorm:"index;not null"`
	Province   Province `gorm:"foreignKey:ProvinceID"`

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

//
// ====== Input ======
//

type DistrictInput struct {
	ID         uint   `json:"id" validate:"required_if=Mode update"`
	Name       string `json:"name" validate:"required"`
	ProvinceID uint   `json:"province_id" validate:"required"`
	Actor      string `json:"-"` // filled in controller
	Mode       string `json:"-"` // "create" or "update"
}

func (input *DistrictInput) Validate() error {
	custom := map[string]string{
		"ID.required_if":      "District ID is required for update",
		"Name.required":       "District name is required",
		"ProvinceID.required": "Province is required",
	}
	return shared.CustomValidate(input, custom)
}

func (input *DistrictInput) ToModel() *District {
	now := time.Now()
	return &District{
		ID:         input.ID,
		Name:       input.Name,
		ProvinceID: input.ProvinceID,
		CreatedBy:  input.Actor,
		UpdatedBy:  input.Actor,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

//
// ====== Update Methods ======
//

func (d *District) UpdateFromInput(input *DistrictInput) {
	d.Name = input.Name
	d.ProvinceID = input.ProvinceID
	d.UpdatedBy = input.Actor
	d.UpdatedAt = time.Now()
}

func (d *District) UpdateFromModel(new *District) {
	d.Name = new.Name
	d.ProvinceID = new.ProvinceID
	d.UpdatedBy = new.UpdatedBy
	d.UpdatedAt = time.Now()
}

func (d *District) MarkDeleted(actor string) {
	d.DeletedBy = actor
	d.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
}

//
// ====== Response ======
//

type DistrictResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ProvinceID   uint   `json:"province_id"`
	ProvinceName string `json:"province_name"`
}

func (d *District) ToResponse() DistrictResponse {
	return DistrictResponse{
		ID:           d.ID,
		Name:         d.Name,
		ProvinceID:   d.ProvinceID,
		ProvinceName: d.Province.Name,
	}
}

func ToDistrictResponses(list []District) []DistrictResponse {
	res := make([]DistrictResponse, len(list))
	for i := range list {
		res[i] = list[i].ToResponse()
	}
	return res
}
