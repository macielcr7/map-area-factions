package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/macielcr7/map-area-factions/backend/internal/services"
)

type HealthHandler struct {
	dbService    *services.DatabaseService
	redisService *services.RedisService
}

type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

func NewHealthHandler(dbService *services.DatabaseService, redisService *services.RedisService) *HealthHandler {
	return &HealthHandler{
		dbService:    dbService,
		redisService: redisService,
	}
}

// Health handles health check endpoint
// @Summary Health check
// @Description Check the health status of the API and its dependencies
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	response := HealthResponse{
		Status:   "healthy",
		Services: make(map[string]string),
	}

	// Check database
	if err := h.dbService.HealthCheck(); err != nil {
		response.Status = "unhealthy"
		response.Services["database"] = "unhealthy"
	} else {
		response.Services["database"] = "healthy"
	}

	// Check Redis
	if err := h.redisService.HealthCheck(); err != nil {
		response.Status = "unhealthy"
		response.Services["redis"] = "unhealthy"
	} else {
		response.Services["redis"] = "healthy"
	}

	// Set appropriate status code
	statusCode := 200
	if response.Status == "unhealthy" {
		statusCode = 503
	}

	return c.Status(statusCode).JSON(response)
}

// LivenessProbe handles Kubernetes liveness probe
// @Summary Liveness probe
// @Description Simple liveness probe for Kubernetes
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health/live [get]
func (h *HealthHandler) LivenessProbe(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "alive",
	})
}

// ReadinessProbe handles Kubernetes readiness probe
// @Summary Readiness probe
// @Description Readiness probe for Kubernetes that checks dependencies
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /health/ready [get]
func (h *HealthHandler) ReadinessProbe(c *fiber.Ctx) error {
	// Check if we can connect to essential services
	if err := h.dbService.HealthCheck(); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status": "not ready",
			"reason": "database connection failed",
		})
	}

	return c.JSON(fiber.Map{
		"status": "ready",
	})
}