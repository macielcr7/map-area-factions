package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/macielcr7/map-area-factions/backend/internal/handlers"
	"github.com/macielcr7/map-area-factions/backend/internal/middleware"
	"github.com/macielcr7/map-area-factions/backend/internal/services"
	"github.com/macielcr7/map-area-factions/backend/internal/utils"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Arquivo .env não encontrado, usando variáveis do sistema")
	}

	// Configurar Viper
	setupConfig()

	// Configurar logs
	setupLogging()

	// Inicializar serviços
	dbService := services.NewDatabaseService()
	if err := dbService.Connect(); err != nil {
		logrus.Fatal("Falha ao conectar com o banco de dados: ", err)
	}
	defer dbService.Close()

	redisService := services.NewRedisService()
	if err := redisService.Connect(); err != nil {
		logrus.Fatal("Falha ao conectar com o Redis: ", err)
	}
	defer redisService.Close()

	// Executar migrações
	if err := dbService.Migrate(); err != nil {
		logrus.Fatal("Falha ao executar migrações: ", err)
	}

	// Configurar Fiber
	app := fiber.New(fiber.Config{
		AppName:      "Map Area Factions API",
		ServerHeader: "Fiber",
		ErrorHandler: middleware.ErrorHandler,
		BodyLimit:    viper.GetInt("max_file_size"),
	})

	// Middleware global
	setupMiddleware(app)

	// Configurar rotas
	setupRoutes(app, dbService, redisService)

	// Iniciar servidor
	port := viper.GetString("port")
	if port == "" {
		port = "8080"
	}

	logrus.Infof("🚀 Servidor iniciando na porta %s", port)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logrus.Info("🛑 Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	// Iniciar servidor
	if err := app.Listen(":" + port); err != nil {
		logrus.Fatal("Falha ao iniciar servidor: ", err)
	}
}

func setupConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	viper.AutomaticEnv()

	// Valores padrão
	viper.SetDefault("port", "8080")
	viper.SetDefault("env", "development")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("jwt_expiry", "24h")
	viper.SetDefault("rate_limit_max", 100)
	viper.SetDefault("rate_limit_window", "1m")
	viper.SetDefault("max_file_size", 10*1024*1024) // 10MB

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warn("Arquivo de configuração não encontrado, usando padrões")
	}
}

func setupLogging() {
	level, err := logrus.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	if viper.GetString("env") == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}

func setupMiddleware(app *fiber.App) {
	// Recovery
	app.Use(recover.New())

	// Logger
	app.Use(logger.New(logger.Config{
		Format: "${time} ${method} ${path} ${status} ${latency}\n",
	}))

	// Security headers
	app.Use(helmet.New())

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     viper.GetString("cors_origins"),
		AllowCredentials: viper.GetBool("cors_credentials"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        viper.GetInt("rate_limit_max"),
		Expiration: utils.ParseDuration(viper.GetString("rate_limit_window")),
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
	}))

	// Prometheus metrics
	if viper.GetBool("prometheus_enabled") {
		app.Use(middleware.PrometheusMiddleware())
	}
}

func setupRoutes(app *fiber.App, db *services.DatabaseService, redis *services.RedisService) {
	// Health check
	app.Get("/health", handlers.HealthCheck(db, redis))
	app.Get("/health/db", handlers.DatabaseHealthCheck(db))
	app.Get("/health/redis", handlers.RedisHealthCheck(redis))

	// Metrics endpoint
	if viper.GetBool("prometheus_enabled") {
		app.Get("/metrics", handlers.PrometheusHandler())
	}

	// API routes
	api := app.Group("/api/v1")

	// Public routes
	api.Post("/auth/login", handlers.Login(db))
	api.Post("/auth/register", handlers.Register(db))
	api.Get("/factions", handlers.GetFactions(db))
	api.Get("/regions", handlers.GetRegions(db))
	api.Get("/search", handlers.Search(db))
	api.Get("/geocode", handlers.Geocode())

	// Protected routes
	protected := api.Use(middleware.JWTMiddleware())
	
	// Auth
	protected.Post("/auth/refresh", handlers.RefreshToken(redis))
	protected.Post("/auth/logout", handlers.Logout(redis))
	protected.Get("/auth/me", handlers.GetProfile(db))

	// User routes
	protected.Get("/geometries", handlers.GetGeometries(db))
	protected.Get("/geometries/:id", handlers.GetGeometry(db))
	protected.Get("/incidents", handlers.GetIncidents(db))
	protected.Get("/incidents/nearby", handlers.GetNearbyIncidents(db))
	protected.Post("/reports", handlers.CreateReport(db))

	// Admin routes
	admin := protected.Use(middleware.AdminMiddleware())
	admin.Post("/factions", handlers.CreateFaction(db))
	admin.Put("/factions/:id", handlers.UpdateFaction(db))
	admin.Delete("/factions/:id", handlers.DeleteFaction(db))
	admin.Post("/geometries", handlers.CreateGeometry(db, redis))
	admin.Put("/geometries/:id", handlers.UpdateGeometry(db, redis))
	admin.Delete("/geometries/:id", handlers.DeleteGeometry(db, redis))
	admin.Post("/geometries/import", handlers.ImportGeoJSON(db, redis))
	admin.Get("/geometries/:id/history", handlers.GetGeometryHistory(db))

	// Subscription routes
	api.Get("/subscriptions/plans", handlers.GetSubscriptionPlans())
	protected.Post("/subscriptions/create", handlers.CreateSubscription(db))
	api.Post("/subscriptions/webhook", handlers.MercadoPagoWebhook(db))

	// WebSocket for real-time updates
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws", websocket.New(handlers.WebSocketHandler(redis)))

	// Server-Sent Events fallback
	app.Get("/events", handlers.SSEHandler(redis))

	// Catch all
	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"error": "Route not found",
		})
	})
}