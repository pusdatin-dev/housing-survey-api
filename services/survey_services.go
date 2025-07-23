package services

import (
	"errors"
	"fmt"
	"housing-survey-api/config"
	"housing-survey-api/internal/context"
	"housing-survey-api/models"
	"housing-survey-api/shared"
	"housing-survey-api/utils"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SurveyService interface {
	GetAllSurveys(ctx *fiber.Ctx) models.ServiceResponse
	GetSurveyDetail(ctx *fiber.Ctx, id string) models.ServiceResponse
	CreateSurvey(ctx *fiber.Ctx, survey models.SurveyInput) models.ServiceResponse
	UpdateSurvey(ctx *fiber.Ctx, survey models.SurveyInput) models.ServiceResponse
	DeleteSurvey(ctx *fiber.Ctx, id string) models.ServiceResponse
	ActionSurvey(ctx *fiber.Ctx, input models.SurveyActionInput) models.ServiceResponse
	GetSurveysByResource(ctx *fiber.Ctx) models.ServiceResponse
	GetSurveysByProgramType(ctx *fiber.Ctx) models.ServiceResponse
	GetSurveysByVerificationStatus(ctx *fiber.Ctx) models.ServiceResponse
}

type surveyService struct {
	Db     *gorm.DB
	Config *config.Config
}

func NewSurveyService(ctx *context.AppContext) SurveyService {
	return &surveyService{
		Db:     ctx.DB,
		Config: ctx.Config,
	}
}

func (s *surveyService) GetAllSurveys(ctx *fiber.Ctx) models.ServiceResponse {
	var surveys []models.Survey
	db := s.Db.Model(&models.Survey{})

	// 1. Ambil role & user id (fallback ke "public" kalau ga login)
	actorRole := "public"
	actorId := uint(0)
	var actor models.User

	// Coba ambil dari JWT context, kalau gagal role public
	if role, err := utils.GetRoleNameFromContext(ctx); err == nil {
		actorRole = role
	}
	if id, err := utils.GetUserIDFromContext(ctx); err == nil {
		actorId = uint(id)
		// Kalau user login, ambil sekalian profile-nya
		if err := s.Db.Preload("Profile").Where("id = ?", actorId).First(&actor).Error; err != nil {
			// biarin, nanti handle di bawah
		}
	}

	// 2. Query role-based filter
	switch actorRole {
	case s.Config.Roles.Surveyor:
		db = db.Where("user_id = ?", actorId)
	case s.Config.Roles.VerificatorBalai, s.Config.Roles.AdminBalai:
		if actor.Profile.ID != 0 {
			db = db.Joins("JOIN profiles ON profiles.user_id = surveys.user_id").
				Where("profiles.balai_id = ?", actor.Profile.BalaiID)
		}
		// case "public", "superadmin", "verificator_eselon1", "admin_eselon1" → akses semua data, tidak perlu filter khusus
	}

	// Filtering
	if address := ctx.Query("address"); address != "" {
		db = db.Where("address LIKE ?", "%"+address+"%")
	}
	if userId := ctx.Query("user_id"); userId != "" {
		db = db.Where("user_id = ?", userId)
	}
	if types := ctx.Query("types"); types != "" {
		// Assuming types is a comma-separated list of survey types
		typeList := utils.SplitAndTrim(types, ",")
		if len(typeList) > 0 {
			db = db.Where("type IN ?", typeList)
		}
	}
	if provinceIDs := ctx.Query("province_ids"); provinceIDs != "" {
		// Assuming province_ids is a comma-separated list of province IDs
		provinceIDList := utils.SplitAndTrim(provinceIDs, ",")
		if len(provinceIDList) > 0 {
			db = db.Where("province_id IN ?", provinceIDList)
		}
	}
	if districtIDs := ctx.Query("district_ids"); districtIDs != "" {
		// Assuming district_ids is a comma-separated list of district IDs
		districtIDList := utils.SplitAndTrim(districtIDs, ",")
		if len(districtIDList) > 0 {
			db = db.Where("district_id IN ?", districtIDList)
		}
	}
	if subdistrictIDs := ctx.Query("subdistrict_ids"); subdistrictIDs != "" {
		// Assuming subdistrict_ids is a comma-separated list of subdistrict IDs
		subdistrictIDList := utils.SplitAndTrim(subdistrictIDs, ",")
		if len(subdistrictIDList) > 0 {
			db = db.Where("subdistrict_id IN ?", subdistrictIDList)
		}
	}
	if villageIDs := ctx.Query("village_ids"); villageIDs != "" {
		// Assuming village_ids is a comma-separated list of village IDs
		villageIDList := utils.SplitAndTrim(villageIDs, ",")
		if len(villageIDList) > 0 {
			db = db.Where("village_id IN ?", villageIDList)
		}
	}
	if programTypeIDs := ctx.Query("program_type_ids"); programTypeIDs != "" {
		// Assuming program_type_ids is a comma-separated list of program type IDs
		programTypeIDList := utils.SplitAndTrim(programTypeIDs, ",")
		if len(programTypeIDList) > 0 {
			db = db.Where("program_type_id IN ?", programTypeIDList)
		}
	}
	if resourceIDs := ctx.Query("resource_ids"); resourceIDs != "" {
		// Assuming resource_ids is a comma-separated list of resource IDs
		resourceIDList := utils.SplitAndTrim(resourceIDs, ",")
		if len(resourceIDList) > 0 {
			db = db.Where("resource_id IN ?", resourceIDList)
		}
	}
	if programIDs := ctx.Query("program_ids"); programIDs != "" {
		// Assuming program_ids is a comma-separated list of program IDs
		programIDList := utils.SplitAndTrim(programIDs, ",")
		if len(programIDList) > 0 {
			db = db.Where("program_id IN ?", programIDList)
		}
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
		return models.InternalServerErrorResponse("Failed to count surveys")
	}
	fmt.Println(db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("User").Limit(limit).Offset(offset).Order("created_at desc").Find(&surveys)
	}))
	if err := db.Preload("User").
		Preload("ProgramType").Preload("Resource").Preload("Program").
		Preload("Province").Preload("District").Preload("Subdistrict").Preload("Village").
		Limit(limit).Offset(offset).Order("created_at desc").
		Find(&surveys).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve surveys")
	}
	// Return with metadata
	return models.OkResponse(fiber.StatusOK, "Survey retrieved successfully", fiber.Map{
		"data":       models.ToSurveyResponse(surveys),
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int((total + int64(limit) - 1) / int64(limit)), // ceiling division
	})
}

func (s *surveyService) GetSurveyDetail(ctx *fiber.Ctx, id string) models.ServiceResponse {
	var survey models.Survey
	if err := s.Db.Preload("User").
		Preload("ProgramType").Preload("Resource").Preload("Program").
		Preload("Province").Preload("District").Preload("Subdistrict").Preload("Village").
		Where("id = ?", id).First(&survey).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve survey")
	}
	if &survey == nil {
		return models.NotFoundResponse("Survey not found")
	}
	return models.OkResponse(fiber.StatusOK, "Survey retrieved successfully", survey.ToResponse())
}

func (s *surveyService) CreateSurvey(ctx *fiber.Ctx, input models.SurveyInput) models.ServiceResponse {
	//enforcing role surveyor only will be in middleware
	// Convert input to model
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}
	survey := input.ToSurvey()

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Cannot find UserID in token")
	}
	utils.LogAudit(ctx, "START", "API entered")
	if userID != int(survey.UserID) {
		return models.BadRequestResponse("Cannot create survey for another user")
	}

	// Insert into DB
	if err := s.Db.Create(&survey).Error; err != nil {
		utils.LogAudit(ctx, "CREATE_SURVEY", err.Error())
		return models.InternalServerErrorResponse("Failed to create survey")
	}

	return models.OkResponse(fiber.StatusCreated, "Survey created successfully", survey.ToResponse())
}

func (s *surveyService) UpdateSurvey(ctx *fiber.Ctx, survey models.SurveyInput) models.ServiceResponse {
	//enforcing role surveyor only will be in middleware
	//newSurvey := survey.ToSurvey()
	oldSurvey := models.Survey{}

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Cannot find UserID in token")
	}

	if userID != int(survey.UserID) {
		return models.BadRequestResponse("Cannot update survey for another user")
	}

	if err := s.Db.Where("id = ? AND deleted_at IS NULL", survey.ID).First(&oldSurvey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Survey not found")
		}
		return models.InternalServerErrorResponse("Failed to retrieve survey for update")
	}

	// Insert into DB
	oldSurvey.UpdateFromInput(survey)
	if err := s.Db.Save(&oldSurvey).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to update survey")
	}

	return models.OkResponse(fiber.StatusCreated, "Survey created successfully", oldSurvey.ToResponse())
}

func (s *surveyService) DeleteSurvey(ctx *fiber.Ctx, id string) models.ServiceResponse {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Cannot find UserID in token")
	}

	var survey models.Survey
	if err = s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&survey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse(fmt.Sprintf("Survey with id %s not found", id))
		}
		return models.InternalServerErrorResponse(fmt.Sprintf("Failed to retrieve survey with id %s", id))
	}

	if userID != int(survey.UserID) {
		return models.BadRequestResponse("Cannot delete survey for another user")
	}

	survey.DeletedBy = fmt.Sprint(userID)
	survey.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
	if err = s.Db.Save(&survey).Error; err != nil {
		return models.InternalServerErrorResponse(fmt.Sprintf("Failed to delete survey with id %s", id))
	}

	return models.OkResponse(200, "Survey deleted successfully", nil)
}

func (s *surveyService) ActionSurvey(ctx *fiber.Ctx, input models.SurveyActionInput) models.ServiceResponse {
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	role, err := utils.GetRoleNameFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Cannot determine role")
	}

	isVerificatorBalai := role == s.Config.Roles.VerificatorBalai
	isVerificatorEselon1 := role == s.Config.Roles.VerificatorEselon1

	if !isVerificatorBalai && !isVerificatorEselon1 {
		return models.ForbiddenResponse("You are not allowed to perform this action")
	}

	// Base query
	db := s.Db.Model(&models.Survey{}).
		Where("id IN ?", input.SurveyIDs).
		Where("is_submitted = ? AND deleted_at IS NULL", true)

	// Filter based on role
	if isVerificatorBalai {
		db = db.Where("status_balai = ?", shared.Pending)
	} else if isVerificatorEselon1 {
		db = db.Where("status_balai = ? AND status_eselon1 = ?", shared.Approved, shared.Pending)
	}

	// Prepare update map
	update := map[string]interface{}{}
	if input.Action == shared.Rejected {
		update["notes"] = input.Notes
	}
	if isVerificatorBalai {
		update["status_balai"] = input.Action
	} else {
		update["status_eselon1"] = input.Action
	}

	// Perform update in one query
	result := db.Updates(update)
	if result.Error != nil {
		return models.InternalServerErrorResponse("Failed to update surveys")
	}

	// Calculate counts
	successCount := result.RowsAffected
	failedCount := int64(len(input.SurveyIDs)) - successCount

	return models.OkResponse(fiber.StatusOK, fmt.Sprintf(
		"%s %d survey(s), %d failed", input.Action, successCount, failedCount,
	), fiber.Map{
		"success_count": successCount,
		"failed_count":  failedCount,
	})
}

func (s *surveyService) GetSurveysByResource(ctx *fiber.Ctx) models.ServiceResponse {
	action := "DASHBOARD_RESOURCE"
	actorRole, err := utils.GetRoleNameFromContext(ctx)
	if err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Cannot get RoleID from context")
	}
	actorId, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Cannot get UserID from context")
	}

	var actor models.User
	if err = s.Db.Preload("Profile").Where("id = ?", actorId).First(&actor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogAudit(ctx, action, err.Error())
			return models.NotFoundResponse("User not found")
		}
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Error retrieving user")
	}

	// 1. Ambil semua resource (buat map tag -> name)
	var resources []models.Resource
	if err = s.Db.Find(&resources).Error; err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Error retrieving resources")
	}
	// Map tag ke name
	tagToName := make(map[string]string)
	for _, r := range resources {
		// hanya isi jika belum ada (biar ambil yang pertama/utama)
		if _, ok := tagToName[r.Tag]; !ok {
			tagToName[r.Tag] = r.Name
		}
	}

	// 2. Hitung survey per tag (bukan per resource_id lagi)
	tagCount := make(map[string]int64)
	for _, r := range resources {
		db := s.Db.Model(&models.Survey{})
		switch actorRole {
		case s.Config.Roles.Surveyor:
			db = db.Where("user_id = ?", actorId)
		case s.Config.Roles.VerificatorBalai, s.Config.Roles.AdminBalai:
			db = db.Joins("JOIN profiles ON profiles.user_id = surveys.user_id").
				Where("profiles.balai_id = ?", actor.Profile.BalaiID)
		}
		var resCount int64
		if err := db.Where("resource_id = ?", r.ID).Count(&resCount).Error; err != nil {
			utils.LogAudit(ctx, action, err.Error())
			return models.InternalServerErrorResponse("cannot count surveys by resource tag")
		}
		tagCount[r.Tag] += resCount // group by tag
	}

	// 3. List tag utama (kalau mau urutan tertentu, bisa manual array ["negara", ...])
	listTag := []string{"negara", "pengembang", "swadaya", "gotongroyong"}

	// 4. Siapkan hasil output (name diambil dari tagToName, total dari tagCount)
	var result []models.DashboardResource
	for _, tag := range listTag {
		result = append(result, models.DashboardResource{
			Name:  tagToName[tag],
			Total: tagCount[tag],
		})
	}

	utils.LogAudit(ctx, action, "Success")
	return models.OkResponse(200, "Success", result)
}

func (s *surveyService) GetSurveysByProgramType(ctx *fiber.Ctx) models.ServiceResponse {
	action := "DASHBOARD_PROGRAM_TYPE"
	actorRole, err := utils.GetRoleNameFromContext(ctx)
	if err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Cannot get RoleID from context")
	}
	actorId, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Cannot get UserID from context")
	}

	var actor models.User
	if err = s.Db.Preload("Profile").Where("id = ?", actorId).First(&actor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogAudit(ctx, action, err.Error())
			return models.NotFoundResponse("User not found")
		}
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Error retrieving user")
	}

	var programTypes []models.ProgramType
	if err = s.Db.Find(&programTypes).Error; err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Error retrieving program types")
	}

	// --- Tambah: Hitung total survey (untuk persentase)
	var totalSurvey int64
	dbAll := s.Db.Model(&models.Survey{})
	switch actorRole {
	case s.Config.Roles.Surveyor:
		dbAll = dbAll.Where("user_id = ?", actorId)
	case s.Config.Roles.VerificatorBalai, s.Config.Roles.AdminBalai:
		dbAll = dbAll.Joins("JOIN profiles ON profiles.user_id = surveys.user_id").
			Where("profiles.balai_id = ?", actor.Profile.BalaiID)
	}
	if err := dbAll.Count(&totalSurvey).Error; err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Error counting total surveys")
	}
	// ---

	var result []models.DashboardProgramType
	for _, pt := range programTypes {
		var resCount int64
		db := s.Db.Model(&models.Survey{})

		switch actorRole {
		case s.Config.Roles.Surveyor:
			db.Where("user_id = ?", actorId)
		case s.Config.Roles.VerificatorBalai, s.Config.Roles.AdminBalai:
			db = db.Joins("JOIN profiles ON profiles.user_id = surveys.user_id").
				Where("profiles.balai_id = ?", actor.Profile.BalaiID)
			//case s.Config.Roles.VerificatorEselon1, s.Config.Roles.AdminEselon1:
		}

		if err = db.Model(&models.Survey{}).Where("program_type_id = ?", pt.ID).Count(&resCount).Error; err != nil {
			utils.LogAudit(ctx, action, err.Error())
			return models.InternalServerErrorResponse("cannot count surveys by resource")
		}

		percent := 0.0
		if totalSurvey > 0 {
			percent = (float64(resCount) / float64(totalSurvey)) * 100
			percent = math.Round(percent*10) / 10 // Satu angka di belakang koma
		}

		result = append(result, models.DashboardProgramType{
			Name:    pt.Name,
			Total:   int(resCount),
			Percent: percent,
		})
	}

	utils.LogAudit(ctx, action, "Success")
	return models.OkResponse(200, "Success", result)
}

func (s *surveyService) GetSurveysByVerificationStatus(ctx *fiber.Ctx) models.ServiceResponse {
	action := "DASHBOARD_VERIFIED"

	// 1. Ambil role dan user_id dari context
	actorRole, err := utils.GetRoleNameFromContext(ctx)
	if err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Cannot get RoleID from context")
	}
	actorId, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Cannot get UserID from context")
	}

	// 2. Ambil data profil actor
	var actor models.User
	if err = s.Db.Preload("Profile").Where("id = ?", actorId).First(&actor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogAudit(ctx, action, err.Error())
			return models.NotFoundResponse("User not found")
		}
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Error retrieving user")
	}

	// 3. Mulai Query dan Filter berdasarkan role actor
	// 3.1 Query untuk mengambil data survey
	db := s.Db.Model(&models.Survey{})

	// 3.2 Filter berdasarkan role actor
	switch actorRole {
	case s.Config.Roles.Surveyor:
		// 3.2.a. Jika Surveyor, hanya survey yang dibuat olehnya
		db.Where("user_id = ?", actorId)
	case s.Config.Roles.VerificatorBalai, s.Config.Roles.AdminBalai:
		// 3.2.b. Jika Verificator Balai atau Admin Balai, tampilkan data Balai
		db.Joins("JOIN profiles ON profiles.user_id = surveys.user_id").
			Where("profiles.balai_id = ?", actor.Profile.BalaiID)
	case s.Config.Roles.VerificatorEselon1, s.Config.Roles.AdminEselon1:
		// 3.2.c. Jika Verificator Eselon1 atau Admin Eselon1, tampilkan semua data
	}

	// 4. Hitung total survey (semua status)
	var totalSurvey int64
	if err := db.Count(&totalSurvey).Error; err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Failed to count all surveys")
	}

	// 5. Hitung yang sudah diverifikasi oleh Eselon 1
	var totalSurveyVerified int64
	dbVerified := db.Session(&gorm.Session{}) // copy semua filter di atas
	if err := dbVerified.Where("status_eselon1 = ?", "Approved").Count(&totalSurveyVerified).Error; err != nil {
		utils.LogAudit(ctx, action, err.Error())
		return models.InternalServerErrorResponse("Failed to count verified surveys")
	}

	// 6. Hitung persentase (dalam persen, dua desimal)
	var percentVerified float64 = 0
	if totalSurvey > 0 {
		percentVerified = (float64(totalSurveyVerified) / float64(totalSurvey)) * 100
		percentVerified = math.Round(percentVerified*10) / 10 // satu angka dibelakang koma
	}

	// 7. Bentuk output JSON
	result := models.DashboardVerified{
		Name:          "Survey Verified Recap",
		Total:         int(totalSurvey),
		VerifiedCount: int(totalSurveyVerified),
		Percent:       percentVerified,
	}

	// 8. Log sukses & return response
	utils.LogAudit(ctx, action, "Success")
	return models.OkResponse(200, "Success", result)
}
