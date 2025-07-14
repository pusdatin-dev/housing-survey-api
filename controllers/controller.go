package controllers

import (
	"housing-survey-api/internal/context"
	"housing-survey-api/services"
)

type ControllerRegistry struct {
	Comment  *CommentController
	AuditLog *AuditLogController
	Survey   *SurveyController
	Auth     *AuthController
	User     *UserController
}

func InitControllers(appCtx *context.AppContext) *ControllerRegistry {
	commentController := &CommentController{Comment: services.NewCommentService(appCtx)}
	auditLogController := &AuditLogController{AuditLog: services.NewAuditLogService(appCtx)}
	surveyController := &SurveyController{Survey: services.NewSurveyService(appCtx)}
	userController := &UserController{User: services.NewUserService(appCtx)}
	authController := &AuthController{Db: appCtx.DB, Config: appCtx.Config}

	return &ControllerRegistry{
		Comment:  commentController,
		AuditLog: auditLogController,
		Survey:   surveyController,
		User:     userController,
		Auth:     authController,
	}
}
