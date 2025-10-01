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
	geometryRepo := repository.NewGeometryRepository(dbService.DB)
	auditRepo := repository.NewAuditRepository(dbService.DB)
	reportRepo := repository.NewReportRepository(dbService.DB)
	settingsRepo := repository.NewSettingsRepository(dbService.DB)
	incidentRepo := repository.NewIncidentRepository(dbService.DB)
	subscriptionRepo := repository.NewSubscriptionRepository(dbService.DB)

	// Initialize WebSocket handler
	wsHandler := handlers.NewWebSocketHandler()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, jwtService)
	userHandler := handlers.NewUserHandler(userRepo)
	factionHandler := handlers.NewFactionHandler(factionRepo)
	geometryHandler := handlers.NewGeometryHandler(geometryRepo, wsHandler)
	auditHandler := handlers.NewAuditHandler(auditRepo)
	reportHandler := handlers.NewReportHandler(reportRepo)
	healthHandler := handlers.NewHealthHandler(dbService, redisService)
	settingsHandler := handlers.NewSettingsHandler(settingsRepo)
	dashboardHandler := handlers.NewDashboardHandler(userRepo, factionRepo, geometryRepo, incidentRepo, subscriptionRepo, reportRepo, auditRepo)

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

	// WebSocket upgrade endpoint
	app.Get("/ws", wsHandler.HandleWebSocket)

	// Health routes (no authentication required)
	health := app.Group("/health")
	health.Get("/", healthHandler.Health)
	health.Get("/live", healthHandler.LivenessProbe)
	health.Get("/ready", healthHandler.ReadinessProbe)

	// API info endpoint
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":        "Map Area Factions API",
			"version":     "1.0.0",
			"description": "API for managing map area factions",
			"endpoints": fiber.Map{
				"health":     "/health",
				"api_v1":     "/api/v1",
				"auth":       "/api/v1/auth",
				"users":      "/api/v1/users",
				"factions":   "/api/v1/factions",
				"geometries": "/api/v1/geometries",
				"audit":      "/api/v1/audit",
				"reports":    "/api/v1/reports",
				"websocket":  "/ws",
				"docs":       "/docs",
			},
			"status": "healthy",
		})
	})

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

	// User routes
	users := api.Group("/users")

	// Protected user routes (admin only)
	usersProtected := users.Group("/")
	usersProtected.Use(middleware.JWTMiddleware(jwtService))
	usersProtected.Use(middleware.RequireRole("admin"))
	usersProtected.Get("/", userHandler.GetUsers)
	usersProtected.Get("/:id", userHandler.GetUser)
	usersProtected.Post("/", userHandler.CreateUser)
	usersProtected.Put("/:id", userHandler.UpdateUser)
	usersProtected.Delete("/:id", userHandler.DeleteUser)

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

	// Geometry routes
	geometries := api.Group("/geometries")

	// Public geometry routes
	geometries.Get("/", geometryHandler.GetGeometries)
	geometries.Get("/:id", geometryHandler.GetGeometry)
	geometries.Get("/search", geometryHandler.SearchGeometries)

	// Protected geometry routes (admin/moderator)
	geometriesProtected := geometries.Group("/")
	geometriesProtected.Use(middleware.JWTMiddleware(jwtService))
	geometriesProtected.Use(middleware.RequireRole("admin", "moderator"))
	geometriesProtected.Post("/", geometryHandler.CreateGeometry)
	geometriesProtected.Put("/:id", geometryHandler.UpdateGeometry)
	geometriesProtected.Delete("/:id", geometryHandler.DeleteGeometry)

	// Audit routes (admin only)
	audit := api.Group("/audit")
	audit.Use(middleware.JWTMiddleware(jwtService))
	audit.Use(middleware.RequireRole("admin"))
	audit.Get("/", auditHandler.GetAuditLogs)
	audit.Get("/:id", auditHandler.GetAuditLog)
	audit.Get("/stats", auditHandler.GetAuditStats)

	// Report routes
	reports := api.Group("/reports")

	// Public report creation (any user or anonymous)
	reports.Post("/", reportHandler.CreateReport)

	// Protected report routes (moderator+)
	reportsProtected := reports.Group("/")
	reportsProtected.Use(middleware.JWTMiddleware(jwtService))
	reportsProtected.Use(middleware.RequireRole("admin", "moderator"))
	reportsProtected.Get("/", reportHandler.GetReports)
	reportsProtected.Get("/:id", reportHandler.GetReport)
	reportsProtected.Put("/:id", reportHandler.UpdateReport)
	reportsProtected.Get("/stats", reportHandler.GetReportStats)

	// Admin-only report deletion
	reportsAdmin := reports.Group("/")
	reportsAdmin.Use(middleware.JWTMiddleware(jwtService))
	reportsAdmin.Use(middleware.RequireRole("admin"))
	reportsAdmin.Delete("/:id", reportHandler.DeleteReport)

	// Settings routes (admin only)
	settings := api.Group("/settings")
	settings.Use(middleware.JWTMiddleware(jwtService))
	settings.Use(middleware.RequireRole("admin"))
	settings.Get("/", settingsHandler.GetSettings)
	settings.Put("/", settingsHandler.UpdateSettings)

	// Dashboard routes (admin only)
	dashboard := api.Group("/dashboard")
	dashboard.Use(middleware.JWTMiddleware(jwtService))
	dashboard.Use(middleware.RequireRole("admin"))
	dashboard.Get("/stats", dashboardHandler.GetDashboardStats)
	dashboard.Get("/activity", dashboardHandler.GetActivityData)
	dashboard.Get("/activities", dashboardHandler.GetRecentActivities)
	dashboard.Get("/alerts", dashboardHandler.GetAlerts)

	// Swagger documentation endpoint
	app.Get("/docs/*", func(c *fiber.Ctx) error {
		return c.SendString("Swagger documentation will be available here")
	})

	// Metrics endpoint for Prometheus
	app.Get("/metrics", func(c *fiber.Ctx) error {
		// Basic metrics in Prometheus format
		metrics := `# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",endpoint="/health"} 1
http_requests_total{method="GET",endpoint="/api"} 1
http_requests_total{method="GET",endpoint="/metrics"} 1

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",endpoint="/health",le="0.1"} 1
http_request_duration_seconds_bucket{method="GET",endpoint="/health",le="0.5"} 1
http_request_duration_seconds_bucket{method="GET",endpoint="/health",le="1.0"} 1
http_request_duration_seconds_bucket{method="GET",endpoint="/health",le="+Inf"} 1
http_request_duration_seconds_sum{method="GET",endpoint="/health"} 0.000205417
http_request_duration_seconds_count{method="GET",endpoint="/health"} 1

# HELP database_connections_active Number of active database connections
# TYPE database_connections_active gauge
database_connections_active 1

# HELP redis_connections_active Number of active Redis connections
# TYPE redis_connections_active gauge
redis_connections_active 1

# HELP application_info Application information
# TYPE application_info gauge
application_info{version="1.0.0",name="map-area-factions"} 1
`
		c.Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		return c.SendString(metrics)
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
