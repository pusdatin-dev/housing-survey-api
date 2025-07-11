package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuditLogRoutes(router fiber.Router, ctrl *controllers.AuditLogController) {
	audit := router.Group("/audit", middleware.AuthRequired(), middleware.AdminOnly())
	audit.Get("/", ctrl.GetAuditLogs)
}
