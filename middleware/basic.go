package middleware

import "github.com/gofiber/fiber/v2"

// Set is a chainable builder for route middleware
type Set struct {
	handlers []fiber.Handler
}

// New initializes a new middleware chain
func New() *Set {
	return &Set{
		handlers: []fiber.Handler{},
	}
}

// Basic adds audit and user injection middlewares
func (m *Set) Basic() *Set {
	m.handlers = append(m.handlers, AuditLogger(), InjectUserAuditFields())
	//m.handlers = append(m.handlers)
	return m
}

// Auth adds authentication + basic middlewares
func (m *Set) Auth() *Set {
	m.handlers = append(m.handlers, AuthRequired())
	return m.Basic()
}

// Public adds public non-auth + basic middlewares
func (m *Set) Public() *Set {
	m.handlers = append(m.handlers, PublicAccess())
	return m.Basic()
}

// Admin adds auth + admin + basic middlewares
func (m *Set) Admin() *Set {
	m.handlers = append(m.handlers, AuthRequired(), AdminOnly())
	return m.Basic()
}

// Custom lets you add arbitrary handlers
func (m *Set) Custom(h ...fiber.Handler) *Set {
	m.handlers = append(m.handlers, h...)
	return m
}

// Build returns the slice of all handlers
func (m *Set) Build() []fiber.Handler {
	return m.handlers
}
