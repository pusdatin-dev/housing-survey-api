package middleware

import (
	"housing-survey-api/internal/context"

	"github.com/gofiber/fiber/v2"
)

var appCtx *context.AppContext

func InitMiddleware(ctx *context.AppContext) {
	appCtx = ctx
}

// Set defines a middleware chain builder
type Set struct {
	handlers []fiber.Handler
}

// New initializes a new middleware chain
func New() *Set {
	return &Set{handlers: []fiber.Handler{}}
}

// Basic appends the audit logger
func (m *Set) Basic() *Set {
	return m.Custom(AuditLogger())
}

// Auth appends AuthRequired and audit logger
func (m *Set) Auth() *Set {
	return m.Custom(AuthRequired()).Basic()
}

// Public appends PublicAccess and audit logger
func (m *Set) Public() *Set {
	return m.Custom(PublicAccess()).Basic()
}

// Admin appends AuthRequired, AdminOnly and audit logger
func (m *Set) Admin() *Set {
	return m.Custom(AuthRequired(), AdminOnly()).Basic()
}

// Custom appends custom handlers
func (m *Set) Custom(h ...fiber.Handler) *Set {
	m.handlers = append(m.handlers, h...)
	return m
}

// Build returns all handlers as variadic slice
func (m *Set) Build() []fiber.Handler {
	return m.handlers
}

// With returns a handler chain in order: middleware..., handler
func With(handler fiber.Handler, middleware ...fiber.Handler) []fiber.Handler {
	return append(middleware, handler)
}

// ↓↓↓ Shorthand helpers below ↓↓↓

// AuthHandler applies auth middleware to the given handler
func AuthHandler(handler fiber.Handler) []fiber.Handler {
	return With(handler, New().Auth().Build()...)
}

// PublicHandler applies public middleware to the given handler
func PublicHandler(handler fiber.Handler) []fiber.Handler {
	return With(handler, New().Public().Build()...)
}

// AdminHandler applies admin middleware to the given handler
func AdminHandler(handler fiber.Handler) []fiber.Handler {
	return With(handler, New().Admin().Build()...)
}

// CustomHandler lets you define an inline chain with a single handler
func CustomHandler(handler fiber.Handler, m ...fiber.Handler) []fiber.Handler {
	return With(handler, m...)
}
