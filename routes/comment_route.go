package routes

import (
	"fmt"
	"housing-survey-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func CommentRoutes(r fiber.Router, ctrl *controllers.CommentController) {
	fmt.Println("Registering comment route with controller:", ctrl != nil)
	comments := r.Group("/comments")

	comments.Get("", ctrl.GetComments)
	comments.Get("/:id", ctrl.GetCommentByID)
	comments.Post("", ctrl.CreatePublicComment)
}
