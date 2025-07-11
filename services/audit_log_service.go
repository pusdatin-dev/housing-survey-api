package services

import (
	"net/http"

	"housing-survey-api/config"
	"housing-survey-api/internal/context"
	"housing-survey-api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuditLogService interface {
	GetAuditLogs(ctx *fiber.Ctx) models.ServiceResponse
}

type auditLogService struct {
	Db     *gorm.DB
	config *config.Config
}

func NewAuditLogService(ctx *context.AppContext) AuditLogService {
	return &auditLogService{
		Db:     ctx.DB,
		config: ctx.Config,
	}
}

func (s *auditLogService) GetAuditLogs(ctx *fiber.Ctx) (res models.ServiceResponse) {
	var logs []models.AuditLog
	res.Status = true
	if err := s.Db.Order("created_at desc").Limit(100).Find(&logs).Error; err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "Failed to retrieve audit logs"
		return res
	}
	res.Code = http.StatusOK
	res.Message = "Audit logs retrieved successfully"
	res.Data = logs
	return
}
