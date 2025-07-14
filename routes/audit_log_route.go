package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuditLogRoutes(r fiber.Router, ctrl *controllers.AuditLogController) {
	r.Get("/audit", middleware.AdminHandler(ctrl.GetAuditLogs)...)
}
