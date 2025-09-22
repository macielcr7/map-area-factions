package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/middleware"
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
	idParam := c.Params("id")
	geometryID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "GEOMETRY_002",
			"message": "Invalid geometry ID",
		})
	}

	geometry, err := h.geometryRepo.GetByID(geometryID)
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

func (h *GeometryHandler) SearchGeometries(c *fiber.Ctx) error {
	latParam := c.Query("lat")
	lngParam := c.Query("lng")
	if latParam == "" || lngParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "GEOMETRY_014",
			"message": "Latitude and longitude are required",
		})
	}

	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation Error",
			"code":    "GEOMETRY_015",
			"message": "Invalid latitude",
		})
	}

	lng, err := strconv.ParseFloat(lngParam, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation Error",
			"code":    "GEOMETRY_016",
			"message": "Invalid longitude",
		})
	}

	radiusParam := c.Query("radius", "1000")
	radius, err := strconv.ParseFloat(radiusParam, 64)
	if err != nil || radius <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation Error",
			"code":    "GEOMETRY_017",
			"message": "Invalid radius",
		})
	}

	searchQuery := c.Query("q")

	geometries, err := h.geometryRepo.SearchNearby(lat, lng, radius, searchQuery)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "GEOMETRY_018",
			"message": "Failed to search geometries",
		})
	}

	return c.JSON(fiber.Map{
		"data": geometries,
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
		GeometryType: req.GeometryType,
		GeoJSON:      req.GeoJSON,
		RiskLevel:    models.RiskLevelLow,
		Status:       models.GeometryStatusDraft,
		Source:       req.Source,
		Notes:        req.Notes,
	}

	regionID, err := parseUUIDPointer(req.RegionID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation Error",
			"code":    "GEOMETRY_005",
			"message": "Invalid region_id",
		})
	}
	geometry.RegionID = regionID

	factionID, err := parseUUIDPointer(req.FactionID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation Error",
			"code":    "GEOMETRY_005",
			"message": "Invalid faction_id",
		})
	}
	geometry.FactionID = factionID

	if req.RiskLevel != nil {
		geometry.RiskLevel = models.RiskLevel(*req.RiskLevel)
	}

	if req.Status != nil {
		geometry.Status = models.GeometryStatus(*req.Status)
	}

	if req.ValidityStart != nil {
		start, err := parseTimePointer(req.ValidityStart)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "Validation Error",
				"code":    "GEOMETRY_005",
				"message": "Invalid validity_start",
			})
		}
		geometry.ValidityStart = start
	}

	if req.ValidityEnd != nil {
		end, err := parseTimePointer(req.ValidityEnd)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "Validation Error",
				"code":    "GEOMETRY_005",
				"message": "Invalid validity_end",
			})
		}
		geometry.ValidityEnd = end
	}

	if authorID := middleware.GetUserID(c); authorID != uuid.Nil {
		geometry.AuthorID = &authorID
	}

	if err := h.geometryRepo.Create(geometry); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "GEOMETRY_006",
			"message": "Failed to create geometry",
		})
	}

	factionBroadcastID := ""
	if geometry.FactionID != nil {
		factionBroadcastID = geometry.FactionID.String()
	}

	h.wsHandler.BroadcastGeometryUpdate(
		geometry.ID.String(),
		factionBroadcastID,
		"created",
	)

	return c.Status(201).JSON(fiber.Map{
		"data": geometry,
	})
}

func (h *GeometryHandler) UpdateGeometry(c *fiber.Ctx) error {
	idParam := c.Params("id")
	geometryID, err := uuid.Parse(idParam)
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

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation Error",
			"code":    "GEOMETRY_005",
			"message": err.Error(),
		})
	}

	geometry, err := h.geometryRepo.GetByID(geometryID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "GEOMETRY_009",
			"message": "Geometry not found",
		})
	}

	if req.RegionID != nil {
		regionID, err := parseUUIDPointer(req.RegionID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "Validation Error",
				"code":    "GEOMETRY_005",
				"message": "Invalid region_id",
			})
		}
		geometry.RegionID = regionID
	}

	if req.GeometryType != nil {
		geometry.GeometryType = *req.GeometryType
	}

	if req.GeoJSON != nil {
		geometry.GeoJSON = *req.GeoJSON
	}

	if req.FactionID != nil {
		factionID, err := parseUUIDPointer(req.FactionID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "Validation Error",
				"code":    "GEOMETRY_005",
				"message": "Invalid faction_id",
			})
		}
		geometry.FactionID = factionID
	}

	if req.RiskLevel != nil {
		geometry.RiskLevel = models.RiskLevel(*req.RiskLevel)
	}

	if req.Status != nil {
		geometry.Status = models.GeometryStatus(*req.Status)
	}

	if req.ValidityStart != nil {
		start, err := parseTimePointer(req.ValidityStart)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "Validation Error",
				"code":    "GEOMETRY_005",
				"message": "Invalid validity_start",
			})
		}
		geometry.ValidityStart = start
	}

	if req.ValidityEnd != nil {
		end, err := parseTimePointer(req.ValidityEnd)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "Validation Error",
				"code":    "GEOMETRY_005",
				"message": "Invalid validity_end",
			})
		}
		geometry.ValidityEnd = end
	}

	if req.Source != nil {
		geometry.Source = *req.Source
	}

	if req.Notes != nil {
		geometry.Notes = *req.Notes
	}

	if err := h.geometryRepo.Update(geometry); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "GEOMETRY_010",
			"message": "Failed to update geometry",
		})
	}

	factionBroadcastID := ""
	if geometry.FactionID != nil {
		factionBroadcastID = geometry.FactionID.String()
	}

	h.wsHandler.BroadcastGeometryUpdate(
		geometry.ID.String(),
		factionBroadcastID,
		"updated",
	)

	return c.JSON(fiber.Map{
		"data": geometry,
	})
}

func (h *GeometryHandler) DeleteGeometry(c *fiber.Ctx) error {
	idParam := c.Params("id")
	geometryID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "GEOMETRY_011",
			"message": "Invalid geometry ID",
		})
	}

	geometry, err := h.geometryRepo.GetByID(geometryID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "GEOMETRY_012",
			"message": "Geometry not found",
		})
	}

	if err := h.geometryRepo.Delete(geometryID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "GEOMETRY_013",
			"message": "Failed to delete geometry",
		})
	}

	factionBroadcastID := ""
	if geometry.FactionID != nil {
		factionBroadcastID = geometry.FactionID.String()
	}

	h.wsHandler.BroadcastGeometryUpdate(
		geometry.ID.String(),
		factionBroadcastID,
		"deleted",
	)

	return c.SendStatus(204)
}

func parseUUIDPointer(value *string) (*uuid.UUID, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	parsed, err := uuid.Parse(*value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func parseTimePointer(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, *value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
