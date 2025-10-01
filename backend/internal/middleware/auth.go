package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/auth"
)

func JWTMiddleware(jwtService *auth.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Authorization header required",
				"code":  "AUTH_001",
			})
		}

		// Extract token from Bearer scheme
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid authorization header format",
				"code":  "AUTH_002",
			})
		}

		tokenString := tokenParts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid token",
				"code":  "AUTH_003",
			})
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func RequireRole(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "User not authenticated",
				"code":  "AUTH_004",
			})
		}

		role, ok := userRole.(string)
		if !ok {
			return c.Status(500).JSON(fiber.Map{
				"error": "Invalid user role type",
				"code":  "AUTH_005",
			})
		}

		// Admin has access to everything
		if role == "admin" {
			return c.Next()
		}

		// Check if user has any of the required roles
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "Insufficient permissions",
			"code":  "AUTH_006",
		})
	}
}

// OptionalJWT middleware allows requests with or without JWT tokens
func OptionalJWT(jwtService *auth.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := tokenParts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Next()
		}

		// Store user info in context if token is valid
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *fiber.Ctx) uuid.UUID {
	userID := c.Locals("userID")
	if userID == nil {
		return uuid.Nil
	}
	return userID.(uuid.UUID)
}

// GetUserRole extracts user role from context
func GetUserRole(c *fiber.Ctx) string {
	role := c.Locals("role")
	if role == nil {
		return ""
	}
	return role.(string)
}