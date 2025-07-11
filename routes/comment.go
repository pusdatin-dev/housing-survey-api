package routes

import (
	"housing-survey-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func CommentRoutes(r fiber.Router, ctrl *controllers.CommentController) {
	comments := r.Group("/comments")

	comments.Get("/", ctrl.GetComments)
	comments.Get("/:id", ctrl.GetCommentByID)
	comments.Post("/", ctrl.CreatePublicComment)
	comments.Put("/:id", ctrl.UpdateComment)
	comments.Delete("/:id", ctrl.DeleteComment)
}
