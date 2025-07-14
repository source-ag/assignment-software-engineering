package middleware

import (
	"net/http"
	"os"
	"strings"

	"meteo-api/internal/models"

	"github.com/gin-gonic/gin"
)

const (
	ReadRole  = "read"
	WriteRole = "write"
)

type AuthMiddleware struct {
	readAPIKey  string
	writeAPIKey string
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		readAPIKey:  getEnv("READ_API_KEY", "read-key-123"),
		writeAPIKey: getEnv("WRITE_API_KEY", "write-key-456"),
	}
}

func (a *AuthMiddleware) RequireAuth(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract API key from "Bearer <key>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "Invalid authorization format. Use 'Bearer <api-key>'",
			})
			c.Abort()
			return
		}

		apiKey := parts[1]
		role := a.validateAPIKey(apiKey)

		if role == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "Invalid API key",
			})
			c.Abort()
			return
		}

		// Check if the role has sufficient permissions
		if !a.hasPermission(role, requiredRole) {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Error:   "Insufficient permissions",
			})
			c.Abort()
			return
		}

		// Store role in context for potential use in handlers
		c.Set("role", role)
		c.Next()
	}
}

func (a *AuthMiddleware) validateAPIKey(apiKey string) string {
	switch apiKey {
	case a.writeAPIKey:
		return WriteRole
	case a.readAPIKey:
		return ReadRole
	default:
		return ""
	}
}

func (a *AuthMiddleware) hasPermission(userRole, requiredRole string) bool {
	// Write role has all permissions
	if userRole == WriteRole {
		return true
	}

	// Read role can only access read endpoints
	if userRole == ReadRole && requiredRole == ReadRole {
		return true
	}

	return false
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}