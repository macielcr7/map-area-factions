package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/macielcr7/map-area-factions/backend/internal/auth"
	"github.com/macielcr7/map-area-factions/backend/internal/config"
	"github.com/macielcr7/map-area-factions/backend/internal/handlers"
	"github.com/macielcr7/map-area-factions/backend/internal/middleware"
	"github.com/macielcr7/map-area-factions/backend/internal/repository"
	"github.com/macielcr7/map-area-factions/backend/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Info("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Configure logger
	logrus.SetLevel(logrus.InfoLevel)
	if cfg.Server.Environment == "development" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	logrus.Info("Starting Map Area Factions API server...")

	// Initialize services
	dbService, err := services.NewDatabaseService(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbService.Close()

	redisService, err := services.NewRedisService(&cfg.Redis)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisService.Close()

	// Auto-migrate database
	if err := dbService.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize JWT service
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiry, cfg.JWT.RefreshExpiry)

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbService.DB)
	factionRepo := repository.NewFactionRepository(dbService.DB)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, jwtService)
	factionHandler := handlers.NewFactionHandler(factionRepo)
	healthHandler := handlers.NewHealthHandler(dbService, redisService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ServerHeader: "Map Factions API",
		AppName:      "Map Area Factions v1.0.0",
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			logrus.WithFields(logrus.Fields{
				"method": c.Method(),
				"path":   c.Path(),
				"ip":     c.IP(),
			}).Error(err)

			return c.Status(code).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"code":    "INTERNAL_001",
				"message": "An unexpected error occurred",
			})
		},
	})

	// Global middleware
	app.Use(middleware.SetupRecovery())
	app.Use(middleware.SetupLogger())
	app.Use(middleware.SetupSecurity())
	app.Use(middleware.SetupCORS())
	app.Use(middleware.SetupRateLimit())

	// Health routes (no authentication required)
	health := app.Group("/health")
	health.Get("/", healthHandler.Health)
	health.Get("/live", healthHandler.LivenessProbe)
	health.Get("/ready", healthHandler.ReadinessProbe)

	// API routes
	api := app.Group("/api/v1")
	api.Use(middleware.SetupAPIRateLimit())

	// Auth routes (no authentication required)
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected auth routes
	authProtected := auth.Group("/")
	authProtected.Use(middleware.JWTMiddleware(jwtService))
	authProtected.Get("/me", authHandler.GetProfile)

	// Faction routes
	factions := api.Group("/factions")
	
	// Public faction routes
	factions.Get("/", factionHandler.GetFactions)
	factions.Get("/:id", factionHandler.GetFaction)

	// Protected faction routes (admin only)
	factionsProtected := factions.Group("/")
	factionsProtected.Use(middleware.JWTMiddleware(jwtService))
	factionsProtected.Use(middleware.RequireRole("admin"))
	factionsProtected.Post("/", factionHandler.CreateFaction)
	factionsProtected.Put("/:id", factionHandler.UpdateFaction)
	factionsProtected.Delete("/:id", factionHandler.DeleteFaction)

	// Swagger documentation endpoint
	app.Get("/docs/*", func(c *fiber.Ctx) error {
		return c.SendString("Swagger documentation will be available here")
	})

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "ROUTE_001",
			"message": "The requested endpoint was not found",
		})
	})

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logrus.Info("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	// Start server
	logrus.Infof("Server starting on port %s", cfg.Server.Port)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	logrus.Info("Server stopped")
}