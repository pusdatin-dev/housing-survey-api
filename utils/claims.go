package utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Claim keys (used in both token and context)
type contextKey string

const (
	UserIDKey   = contextKey("user_id")
	EmailKey    = contextKey("user_email")
	NameKey     = contextKey("user_name")
	RoleIDKey   = contextKey("role_id")
	RoleNameKey = contextKey("role_name")
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// ======================
// Context Setters
// ======================

// SetUserContext stores user info from claims into the request context
func SetUserContext(ctx context.Context, claims jwt.MapClaims) context.Context {
	ctx = context.WithValue(ctx, UserIDKey, claims["user_id"])
	ctx = context.WithValue(ctx, EmailKey, claims["user_email"])
	ctx = context.WithValue(ctx, NameKey, claims["user_name"])
	ctx = context.WithValue(ctx, RoleIDKey, claims["role_id"])
	ctx = context.WithValue(ctx, RoleNameKey, claims["role_name"])
	return ctx
}

// SetGuestContext stores guest context (e.g. for public routes)
func SetGuestContext(ctx context.Context, ip string) context.Context {
	ctx = context.WithValue(ctx, UserIDKey, "00000000-0000-0000-0000-000000000000")
	ctx = context.WithValue(ctx, EmailKey, "public")
	ctx = context.WithValue(ctx, NameKey, ip)
	ctx = context.WithValue(ctx, RoleIDKey, "0")
	ctx = context.WithValue(ctx, RoleNameKey, "Guest")
	return ctx
}

// ======================
// JWT Claims Extraction
// ======================

// GetClaims parses JWT claims from Authorization header
func GetClaims(c *fiber.Ctx) (jwt.MapClaims, error) {
	tokenStr, err := extractBearerToken(c)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unable to parse token claims")
	}

	return claims, nil
}

// extractBearerToken pulls Bearer token from header
func extractBearerToken(c *fiber.Ctx) (string, error) {
	auth := c.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		return "", fmt.Errorf("missing or invalid Authorization header")
	}
	return strings.TrimPrefix(auth, "Bearer "), nil
}

// ======================
// Context Access Helpers
// ======================

func GetUserIDFromContext(c *fiber.Ctx) (string, error) {
	return getClaimString(c, "user_id")
}

func GetUserEmailFromContext(c *fiber.Ctx) (string, error) {
	return getClaimString(c, "user_email")
}

func GetUserNameFromContext(c *fiber.Ctx) (string, error) {
	return getClaimString(c, "user_name")
}

func GetRoleIDFromContext(c *fiber.Ctx) (string, error) {
	return getClaimString(c, "role_id")
}

func GetRoleNameFromContext(c *fiber.Ctx) (string, error) {
	return getClaimString(c, "role_name")
}

// getClaimString is a shared internal helper
func getClaimString(c *fiber.Ctx, key string) (string, error) {
	claims, err := GetClaims(c)
	if err != nil {
		return "", err
	}
	val, ok := claims[key]
	if !ok {
		return "", fmt.Errorf("claim '%s' not found", key)
	}
	return fmt.Sprint(val), nil
}

// ======================
// Auth Check Helper
// ======================

// IsAuthenticated returns true if a valid token is present
func IsAuthenticated(c *fiber.Ctx) bool {
	_, err := GetClaims(c)
	return err == nil
}
