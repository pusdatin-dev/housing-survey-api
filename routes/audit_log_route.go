package routes

import (
	"housing-survey-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuditLogRoutes(router fiber.Router, ctrl *controllers.AuditLogController) {
	//audit := router.Group("/audit", middleware.New().Admin().Build()...)
	audit := router.Group("/audit")
	audit.Get("/", ctrl.GetAuditLogs)
}
