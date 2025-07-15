package routes

import (
	"fmt"
	"housing-survey-api/controllers"
	"reflect"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, ctrl *controllers.ControllerRegistry) {
	v1 := app.Group("/api/v1")

	AuthRoutes(v1, ctrl.Auth) // /login, /signup
	UserRoutesV1(v1, ctrl.User)
	CommentRoutes(v1, ctrl.Comment)
	SurveyRoutesV1(v1, ctrl.Survey)
	AuditLogRoutes(v1, ctrl.AuditLog)
	SurveyorRoutes(v1, ctrl.Surveyor)
}

func PrintRoutes(app *fiber.App) {
	for _, route := range app.GetRoutes() {
		fmt.Printf("Route: %-7s %-31s | Handlers: %d\n", route.Method, route.Path, len(route.Handlers))

		for i, handler := range route.Handlers {
			handlerName := getFunctionName(handler)
			fmt.Printf("   └── Handler %d: %s\n", i+1, handlerName)
		}
	}
}

func getFunctionName(f fiber.Handler) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
