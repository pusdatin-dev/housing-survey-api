package routes

import (
	"housing-survey-api/controllers"
	"housing-survey-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func CommentRoutes(r fiber.Router, ctrl *controllers.CommentController) {
	comments := r.Group("/comments")

	comments.Get("", middleware.PublicHandler(ctrl.GetComments)...)
	comments.Get("/:id", middleware.PublicHandler(ctrl.GetCommentByID)...)
	comments.Post("", middleware.PublicHandler(ctrl.CreatePublicComment)...)

	// --> add authenticated routes for verificator roles
	// --> can comment to publicComments and update comments

	//comments.Get("", ctrl.GetComments)
	//comments.Get("/:id", ctrl.GetCommentByID)
	//comments.Post("", ctrl.CreatePublicComment)
}
