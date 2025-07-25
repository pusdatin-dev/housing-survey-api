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
	action := "GET ALL SURVEYS"
	var surveys []models.Survey
	db := s.Db.Model(&models.Survey{})

	// 🔐 Ambil role & user ID dari JWT (fallback ke public jika gagal)
	actorRole := "public"
	actorId := uint(0)
	var actor models.User

	// 🔍 Ambil role user
	if role, err := utils.GetRoleNameFromContext(ctx); err == nil {
		actorRole = role
	}
	// 🔍 Ambil id user
	if id, err := utils.GetUserIDFromContext(ctx); err == nil {
		actorId = uint(id)
		// Kalau user login, ambil sekalian profile-nya
		if err := s.Db.Preload("Profile").Where("id = ?", actorId).First(&actor).Error; err != nil {
		}
	}

	// 🔐 Role-based filtering
	switch actorRole {
	case s.Config.Roles.Surveyor:
		// 👤 Surveyor hanya bisa lihat survei miliknya
		db = db.Where("user_id = ?", actorId)
	case s.Config.Roles.VerificatorBalai, s.Config.Roles.AdminBalai:
		// 🏢 Verifikator/Admin Balai hanya lihat survei di balainya
		if actor.Profile.ID != 0 {
			db = db.Joins("JOIN profiles ON profiles.user_id = surveys.user_id").
				Where("profiles.balai_id = ?", actor.Profile.BalaiID)
		}
	case s.Config.Roles.SuperAdmin:
		// 👑 SuperAdmin bebas akses semua data — tanpa filter
		// Do nothing
	}

	// 🔍 Filter dinamis dari query string
	if address := ctx.Query("address"); address != "" {
		db = db.Where("address LIKE ?", "%"+address+"%")
	}
	// 🚫 Filter user_id hanya untuk role non-surveyor (biar gak bisa iseng inject)
	if userId := ctx.Query("user_id"); userId != "" {
		if actorRole == s.Config.Roles.AdminEselon1 || actorRole == s.Config.Roles.SuperAdmin {
			db = db.Where("user_id = ?", userId)
		}
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

	// 📄 Pagination
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// 🔢 Hitung total data
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to count surveys")
	}

	// 📦 Ambil data dengan preload relasi
	if err := db.Preload("User").
		Preload("ProgramType").Preload("Resource").Preload("Program").
		Preload("Province").Preload("District").Preload("Subdistrict").Preload("Village").
		Limit(limit).Offset(offset).Order("created_at desc").
		Find(&surveys).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to retrieve surveys")
	}

	// 📝 Audit log
	utils.LogAudit(ctx, action, "Success")

	// ✅ Return data + metadata paginasi
	return models.OkResponse(fiber.StatusOK, "Survey retrieved successfully", fiber.Map{
		"data":       models.ToSurveyResponse(surveys),
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int((total + int64(limit) - 1) / int64(limit)), // ceiling division
		"roleid":     actorRole,
		"userid":     actorId,
	})
}

func (s *surveyService) GetSurveyDetail(ctx *fiber.Ctx, id string) models.ServiceResponse {
	action := "GET_SURVEY_DETAIL"
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
	// 8. Log sukses & return response
	utils.LogAudit(ctx, action, "Success")
	return models.OkResponse(fiber.StatusOK, "Survey retrieved successfully", survey.ToResponse())
}

func (s *surveyService) CreateSurvey(ctx *fiber.Ctx, input models.SurveyInput) models.ServiceResponse {
	action := "CREATE_SURVEY"
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
	// 8. Log sukses & return response
	utils.LogAudit(ctx, action, "Success")
	return models.OkResponse(fiber.StatusCreated, "Survey created successfully", survey.ToResponse())
}

func (s *surveyService) UpdateSurvey(ctx *fiber.Ctx, survey models.SurveyInput) models.ServiceResponse {
	action := "UPDATE_SURVEY"

	// 🔐 Ambil user ID dari token JWT
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Cannot find UserID in token")
	}

	// 🚫 Cegah user mengubah survei milik user lain
	if userID != int(survey.UserID) {
		return models.BadRequestResponse("Cannot update survey for another user")
	}

	// 🔍 Ambil data survei lama dari database (yang belum dihapus)
	oldSurvey := models.Survey{}

	if err := s.Db.Where("id = ? AND deleted_at IS NULL", survey.ID).First(&oldSurvey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse("Survey not found")
		}
		return models.InternalServerErrorResponse("Failed to retrieve survey for update")
	}

	// 🚫 Cek apakah survey ditolak oleh Balai
	if oldSurvey.StatusBalai == shared.Rejected {
		return models.ForbiddenResponse("Rejected survey can not be updated")
	}

	// 🔁 Mapping input ke model survei (hanya field yang boleh diubah)
	oldSurvey.UpdateFromInput(survey)

	// ✍️ Tambahkan informasi audit: siapa yang update & kapan
	oldSurvey.UpdatedBy = fmt.Sprint(userID)
	oldSurvey.UpdatedAt = time.Now() // gunakan pointer karena tipe UpdatedAt-nya *time.Time

	// 💾 Simpan perubahan ke database (gunakan .Updates agar aman)
	if err := s.Db.Save(&oldSurvey).Error; err != nil {
		return models.InternalServerErrorResponse("Failed to update survey")
	}

	// 📝 Catat ke audit log
	utils.LogAudit(ctx, action, "Success")

	// ✅ Berhasil: kembalikan respon OK dan data survei terbaru
	return models.OkResponse(fiber.StatusCreated, "Survey created successfully", oldSurvey.ToResponse())
}

func (s *surveyService) DeleteSurvey(ctx *fiber.Ctx, id string) models.ServiceResponse {
	action := "DELETE_SURVEY"

	// 🔐 Ambil user ID dari token JWT
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Cannot find UserID in token")
	}

	// 🔍 Ambil survei berdasarkan ID, pastikan belum dihapus
	var survey models.Survey
	if err = s.Db.Where("id = ? AND deleted_at IS NULL", id).First(&survey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotFoundResponse(fmt.Sprintf("Survey with id %s not found", id))
		}
		return models.InternalServerErrorResponse(fmt.Sprintf("Failed to retrieve survey with id %s", id))
	}

	// 🚫 Cegah user menghapus survei milik orang lain
	if userID != int(survey.UserID) {
		return models.BadRequestResponse("Cannot delete survey for another user")
	}

	// 🔒 Cegah penghapusan jika survei sudah disubmit
	if survey.IsSubmitted {
		return models.ForbiddenResponse("Survei sudah disubmit dan tidak dapat dihapus")
	}

	// 🗑️ Tandai survei sebagai dihapus (soft delete)
	survey.DeletedBy = fmt.Sprint(userID)
	survey.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
	if err = s.Db.Save(&survey).Error; err != nil {
		return models.InternalServerErrorResponse(fmt.Sprintf("Failed to delete survey with id %s", id))
	}

	// 📝 Catat aksi penghapusan ke audit log
	utils.LogAudit(ctx, action, "Success")

	// ✅ Berhasil
	return models.OkResponse(200, "Survey deleted successfully", nil)
}

func (s *surveyService) ActionSurvey(ctx *fiber.Ctx, input models.SurveyActionInput) models.ServiceResponse {
	action := "ACTION_SURVEY"

	// ✅ Validasi input dari request body
	if err := input.Validate(); err != nil {
		return models.BadRequestResponse(err.Error())
	}

	// 🔐 Ambil role & user ID dari JWT
	role, err := utils.GetRoleNameFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Gagal menentukan role user")
	}
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return models.InternalServerErrorResponse("Gagal mengambil User ID dari token")
	}

	isVerificatorBalai := role == s.Config.Roles.VerificatorBalai
	isVerificatorEselon1 := role == s.Config.Roles.VerificatorEselon1

	// 🚫 Cegah jika user bukan verifikator
	if !isVerificatorBalai && !isVerificatorEselon1 {
		return models.ForbiddenResponse("Kamu tidak memiliki izin untuk melakukan aksi ini")
	}

	// 🏢 Ambil balai_id verifikator (jika role adalah verifikator balai)
	var allowedBalaiID uint
	if isVerificatorBalai {
		var profile models.Profile
		if err := s.Db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
			return models.InternalServerErrorResponse("Gagal mengambil profil user")
		}
		if profile.BalaiID == nil {
			return models.ForbiddenResponse("User belum terdaftar ke dalam balai manapun")
		}
		allowedBalaiID = *profile.BalaiID
	}

	// 🔍 Ambil semua survei yang akan diproses, preload relasi user dan profil
	var surveys []models.Survey
	if err := s.Db.Preload("User.Profile").
		Where("id IN ?", input.SurveyIDs).
		Where("is_submitted = ? AND deleted_at IS NULL", true).
		Find(&surveys).Error; err != nil {
		return models.InternalServerErrorResponse("Gagal mengambil data survei")
	}

	successCount := int64(0)
	failedDetails := []string{}

	for _, survey := range surveys {
		// 📋 Validasi status & balai berdasarkan role
		if isVerificatorBalai {
			if survey.StatusBalai != shared.Pending {
				failedDetails = append(failedDetails, fmt.Sprintf("Survey %d: status balai bukan pending", survey.ID))
				continue
			}
			if survey.User.Profile.BalaiID == nil || *survey.User.Profile.BalaiID != allowedBalaiID {
				failedDetails = append(failedDetails, fmt.Sprintf("Survey %d: bukan milik balai anda", survey.ID))
				continue
			}
		} else if isVerificatorEselon1 {
			if !(survey.StatusBalai == shared.Approved && survey.StatusEselon1 == shared.Pending) {
				failedDetails = append(failedDetails, fmt.Sprintf("Survey %d: belum disetujui balai atau sudah diverifikasi", survey.ID))
				continue
			}
		}

		// 🛠️ Siapkan field yang akan diupdate
		update := map[string]interface{}{
			"updated_by": fmt.Sprint(userID), // 👤 ID verifikator
			"updated_at": time.Now(),         // ⏰ waktu update
		}
		if input.Action == shared.Rejected {
			update["notes"] = input.Notes // 📝 alasan ditolak
		}
		if isVerificatorBalai {
			update["status_balai"] = input.Action
		} else {
			update["status_eselon1"] = input.Action
		}

		// 💾 Simpan perubahan
		if err := s.Db.Model(&models.Survey{}).Where("id = ?", survey.ID).Updates(update).Error; err != nil {
			failedDetails = append(failedDetails, fmt.Sprintf("Survey %d: gagal update", survey.ID))
			continue
		}

		successCount++
	}

	// 📝 Catat aksi ke audit log
	utils.LogAudit(ctx, action, "Success")

	// 📦 Kirim respon ke client
	return models.OkResponse(fiber.StatusOK, fmt.Sprintf(
		"%s %d survei, %d gagal", input.Action, successCount, len(failedDetails),
	), fiber.Map{
		"success_count":  successCount,
		"failed_count":   len(failedDetails),
		"failed_details": failedDetails,
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
