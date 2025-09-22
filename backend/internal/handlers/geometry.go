package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"github.com/macielcr7/map-area-factions/backend/internal/repository"
	"github.com/macielcr7/map-area-factions/backend/internal/utils"
)

type GeometryHandler struct {
	geometryRepo *repository.GeometryRepository
	wsHandler    *WebSocketHandler
}

func NewGeometryHandler(geometryRepo *repository.GeometryRepository, wsHandler *WebSocketHandler) *GeometryHandler {
	return &GeometryHandler{
		geometryRepo: geometryRepo,
		wsHandler:    wsHandler,
	}
}

func (h *GeometryHandler) GetGeometries(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	factionID := c.Query("faction_id")

	geometries, total, err := h.geometryRepo.GetAll(page, limit, factionID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "GEOMETRY_001",
			"message": "Failed to fetch geometries",
		})
	}

	return c.JSON(fiber.Map{
		"data": geometries,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *GeometryHandler) GetGeometry(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "GEOMETRY_002",
			"message": "Invalid geometry ID",
		})
	}

	geometry, err := h.geometryRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "GEOMETRY_003",
			"message": "Geometry not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": geometry,
	})
}

func (h *GeometryHandler) CreateGeometry(c *fiber.Ctx) error {
	var req models.CreateGeometryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "GEOMETRY_004",
			"message": "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation Error",
			"code":    "GEOMETRY_005",
			"message": err.Error(),
		})
	}

	geometry := &models.Geometry{
		FactionID:   req.FactionID,
		Name:        req.Name,
		Description: req.Description,
		GeoJSON:     req.GeoJSON,
		Color:       req.Color,
		IsActive:    true,
	}

	if err := h.geometryRepo.Create(geometry); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "GEOMETRY_006",
			"message": "Failed to create geometry",
		})
	}

	// Broadcast real-time update
	h.wsHandler.BroadcastGeometryUpdate(
		strconv.FormatUint(uint64(geometry.ID), 10),
		strconv.FormatUint(uint64(geometry.FactionID), 10),
		"created",
	)

	return c.Status(201).JSON(fiber.Map{
		"data": geometry,
	})
}

func (h *GeometryHandler) UpdateGeometry(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "GEOMETRY_007",
			"message": "Invalid geometry ID",
		})
	}

	var req models.UpdateGeometryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "GEOMETRY_008",
			"message": "Invalid request body",
		})
	}

	geometry, err := h.geometryRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "GEOMETRY_009",
			"message": "Geometry not found",
		})
	}

	// Update fields
	if req.Name != nil {
		geometry.Name = *req.Name
	}
	if req.Description != nil {
		geometry.Description = req.Description
	}
	if req.GeoJSON != nil {
		geometry.GeoJSON = *req.GeoJSON
	}
	if req.Color != nil {
		geometry.Color = req.Color
	}
	if req.IsActive != nil {
		geometry.IsActive = *req.IsActive
	}

	if err := h.geometryRepo.Update(geometry); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "GEOMETRY_010",
			"message": "Failed to update geometry",
		})
	}

	// Broadcast real-time update
	h.wsHandler.BroadcastGeometryUpdate(
		strconv.FormatUint(uint64(geometry.ID), 10),
		strconv.FormatUint(uint64(geometry.FactionID), 10),
		"updated",
	)

	return c.JSON(fiber.Map{
		"data": geometry,
	})
}

func (h *GeometryHandler) DeleteGeometry(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "GEOMETRY_011",
			"message": "Invalid geometry ID",
		})
	}

	geometry, err := h.geometryRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "GEOMETRY_012",
			"message": "Geometry not found",
		})
	}

	if err := h.geometryRepo.Delete(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "GEOMETRY_013",
			"message": "Failed to delete geometry",
		})
	}

	// Broadcast real-time update
	h.wsHandler.BroadcastGeometryUpdate(
		strconv.FormatUint(uint64(geometry.ID), 10),
		strconv.FormatUint(uint64(geometry.FactionID), 10),
		"deleted",
	)

	return c.SendStatus(204)
}

func (h *GeometryHandler) SearchGeometries(c *fiber.Ctx) error {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	radius, _ := strconv.ParseFloat(c.Query("radius", "1000"), 64) // Default 1km
	query := c.Query("q")

	geometries, err := h.geometryRepo.SearchNearby(lat, lng, radius, query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "GEOMETRY_014",
			"message": "Failed to search geometries",
		})
	}

	return c.JSON(fiber.Map{
		"data": geometries,
	})
}