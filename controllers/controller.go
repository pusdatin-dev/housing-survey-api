package controllers

import (
	"housing-survey-api/internal/context"
	"housing-survey-api/services"
)

type ControllerRegistry struct {
	Comment     *CommentController
	AuditLog    *AuditLogController
	Survey      *SurveyController
	Auth        *AuthController
	User        *UserController
	Balai       *BalaiController
	District    *DistrictController
	Program     *ProgramController
	ProgramType *ProgramTypeController
	Province    *ProvinceController
	Resource    *ResourceController
	Role        *RoleController
	Subdistrict *SubdistrictController
	Village     *VillageController
}

func InitControllers(appCtx *context.AppContext) *ControllerRegistry {
	return &ControllerRegistry{
		Comment:     &CommentController{Comment: services.NewCommentService(appCtx)},
		AuditLog:    &AuditLogController{AuditLog: services.NewAuditLogService(appCtx)},
		Survey:      &SurveyController{Survey: services.NewSurveyService(appCtx)},
		Auth:        &AuthController{Db: appCtx.DB, Config: appCtx.Config},
		User:        &UserController{User: services.NewUserService(appCtx)},
		Balai:       &BalaiController{Service: services.NewBalaiService(appCtx)},
		District:    &DistrictController{Service: services.NewDistrictService(appCtx)},
		Program:     &ProgramController{Service: services.NewProgramService(appCtx)},
		ProgramType: &ProgramTypeController{Service: services.NewProgramTypeService(appCtx)},
		Province:    &ProvinceController{Service: services.NewProvinceService(appCtx)},
		Resource:    &ResourceController{Service: services.NewResourceService(appCtx)},
		Role:        &RoleController{Service: services.NewRoleService(appCtx)},
		Subdistrict: &SubdistrictController{Service: services.NewSubdistrictService(appCtx)},
		Village:     &VillageController{Service: services.NewVillageService(appCtx)},
	}
}
