package models

import (
	"time"

	"housing-survey-api/shared"

	"gorm.io/gorm"
)

// Village Master Data
// Village = kelurahan
type Village struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	Name          string `gorm:"type:text;index"`
	SubdistrictID uint   `gorm:"index"`
	Subdistrict   Subdistrict

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (v *Village) UpdateFromInput(input *VillageInput) {
	v.Name = input.Name
	v.SubdistrictID = input.SubdistrictID
	v.UpdatedBy = input.Actor
	v.UpdatedAt = time.Now()
}

func (v *Village) Update(newVillage *Village) {
	v.Name = newVillage.Name
	v.SubdistrictID = newVillage.SubdistrictID
	v.UpdatedBy = newVillage.UpdatedBy
	v.UpdatedAt = time.Now()
}

func (v *Village) MarkDeleted(actor string) {
	v.DeletedBy = actor
	v.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
}

func (v *Village) ToResponse() VillageResponse {
	return VillageResponse{
		ID:              v.ID,
		Name:            v.Name,
		SubdistrictID:   v.SubdistrictID,
		SubdistrictName: v.Subdistrict.Name,
	}
}

func ToVillageResponses(villages []Village) []VillageResponse {
	res := make([]VillageResponse, len(villages))
	for i, v := range villages {
		res[i] = v.ToResponse()
	}
	return res
}

type VillageResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	SubdistrictID   uint   `json:"subdistrict_id"`
	SubdistrictName string `json:"subdistrict_name"`
}

type VillageInput struct {
	ID            uint   `json:"id"`
	Name          string `json:"name" validate:"required"`
	SubdistrictID uint   `json:"subdistrict_id" validate:"required"`
	Actor         string `json:"-"`
	Mode          string `json:"-"`
}

func (v *VillageInput) Validate() error {
	custom := map[string]string{
		"Name.required":          "Village name is required",
		"SubdistrictID.required": "Subdistrict is required",
	}
	return shared.CustomValidate(v, custom)
}

func (v *VillageInput) ToModel() Village {
	now := time.Now()
	return Village{
		ID:            v.ID,
		Name:          v.Name,
		SubdistrictID: v.SubdistrictID,
		CreatedBy:     v.Actor,
		UpdatedBy:     v.Actor,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
