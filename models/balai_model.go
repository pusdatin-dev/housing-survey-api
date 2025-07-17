package models

import (
	"time"

	"housing-survey-api/shared"

	"gorm.io/gorm"
)

type Balai struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	Name          string `gorm:"not null"`
	ProvinceID    uint   `gorm:"index"`
	DistrictID    uint   `gorm:"index"`
	SubdistrictID uint   `gorm:"index"`
	VillageID     uint   `gorm:"index"`
	Province      Province
	District      District
	Subdistrict   Subdistrict
	Village       Village

	CreatedBy string `gorm:"type:text"`
	UpdatedBy string `gorm:"type:text"`
	DeletedBy string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *Balai) UpdateFromInput(input *BalaiInput) {
	b.Name = input.Name
	b.ProvinceID = input.ProvinceID
	b.DistrictID = input.DistrictID
	b.SubdistrictID = input.SubdistrictID
	b.VillageID = input.VillageID
	b.UpdatedBy = input.Actor
	b.UpdatedAt = time.Now()
}

func (b *Balai) UpdateFromModel(newBalai *Balai) {
	b.Name = newBalai.Name
	b.ProvinceID = newBalai.ProvinceID
	b.DistrictID = newBalai.DistrictID
	b.SubdistrictID = newBalai.SubdistrictID
	b.VillageID = newBalai.VillageID
	b.UpdatedBy = newBalai.UpdatedBy
	b.UpdatedAt = time.Now()
}

func (b *Balai) MarkDeleted(actor string) {
	b.DeletedBy = actor
	b.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
}

func (b *Balai) ToResponse() BalaiResponse {
	return BalaiResponse{
		ID:              b.ID,
		Name:            b.Name,
		ProvinceID:      b.ProvinceID,
		ProvinceName:    b.Province.Name,
		DistrictID:      b.DistrictID,
		DistrictName:    b.District.Name,
		SubdistrictID:   b.SubdistrictID,
		SubdistrictName: b.Subdistrict.Name,
		VillageID:       b.VillageID,
		VillageName:     b.Village.Name,
	}
}

func ToBalaiResponses(balais []Balai) []BalaiResponse {
	res := make([]BalaiResponse, len(balais))
	for i, b := range balais {
		res[i] = b.ToResponse()
	}
	return res
}

type BalaiResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	ProvinceID      uint   `json:"province_id"`
	ProvinceName    string `json:"province_name"`
	DistrictID      uint   `json:"district_id"`
	DistrictName    string `json:"district_name"`
	SubdistrictID   uint   `json:"subdistrict_id"`
	SubdistrictName string `json:"subdistrict_name"`
	VillageID       uint   `json:"village_id"`
	VillageName     string `json:"village_name"`
}

type BalaiInput struct {
	ID            uint   `json:"id"`
	Name          string `json:"name" validate:"required"`
	ProvinceID    uint   `json:"province_id" validate:"required"`
	DistrictID    uint   `json:"district_id" validate:"required"`
	SubdistrictID uint   `json:"subdistrict_id" validate:"required"`
	VillageID     uint   `json:"village_id" validate:"required"`
	Actor         string `json:"-"` // CreatedBy / UpdatedBy
	Mode          string `json:"-"` // "create" / "update"
}

func (b *BalaiInput) Validate() error {
	customMessages := map[string]string{
		"Name.required":          "Balai name is required",
		"ProvinceID.required":    "Province is required",
		"DistrictID.required":    "District is required",
		"SubdistrictID.required": "Subdistrict is required",
		"VillageID.required":     "Village is required",
	}
	return shared.CustomValidate(b, customMessages)
}

func (b *BalaiInput) ToModel() Balai {
	now := time.Now()
	return Balai{
		ID:            b.ID,
		Name:          b.Name,
		ProvinceID:    b.ProvinceID,
		DistrictID:    b.DistrictID,
		SubdistrictID: b.SubdistrictID,
		VillageID:     b.VillageID,
		CreatedBy:     b.Actor,
		UpdatedBy:     b.Actor,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
