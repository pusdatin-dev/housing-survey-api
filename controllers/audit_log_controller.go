package controllers

import (
	"housing-survey-api/services"
	"housing-survey-api/utils"

	"github.com/gofiber/fiber/v2"
)

type AuditLogController struct {
	AuditLog services.AuditLogService
}

func (c *AuditLogController) GetAuditLogs(ctx *fiber.Ctx) error {
	return utils.ToFiberJSON(ctx, c.AuditLog.GetAuditLogs(ctx))
}
