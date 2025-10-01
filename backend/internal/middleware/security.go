package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupCORS configures CORS middleware
func SetupCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,http://localhost:3001,https://admin.mapfactions.com",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	})
}

// SetupSecurity configures security headers
func SetupSecurity() fiber.Handler {
	return helmet.New(helmet.Config{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		CrossOriginEmbedderPolicy: "require-corp",
	})
}

// SetupRateLimit configures rate limiting
func SetupRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			// Skip rate limiting for health checks
			return c.Path() == "/health"
		},
		Max:        100, // 100 requests
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Rate limit exceeded",
				"code":  "RATE_001",
			})
		},
	})
}

// SetupAPIRateLimit configures stricter rate limiting for API routes
func SetupAPIRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        50, // 50 requests
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use user ID if authenticated, otherwise IP
			userID := c.Locals("userID")
			if userID != nil {
				return userID.(string)
			}
			return c.Get("x-forwarded-for", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "API rate limit exceeded",
				"code":  "RATE_002",
			})
		},
	})
}

// SetupLogger configures request logging
func SetupLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} ${method} ${path} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "America/Sao_Paulo",
	})
}

// SetupRecovery configures panic recovery
func SetupRecovery() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}