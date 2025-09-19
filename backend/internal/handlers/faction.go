package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"github.com/macielcr7/map-area-factions/backend/internal/repository"
	"github.com/macielcr7/map-area-factions/backend/internal/utils"
)

type FactionHandler struct {
	factionRepo repository.FactionRepository
}

type CreateFactionRequest struct {
	Name            string `json:"name" validate:"required,min=2,max=100"`
	Acronym         string `json:"acronym" validate:"required,min=2,max=10"`
	ColorHex        string `json:"color_hex" validate:"required,hexcolor"`
	DisplayPriority int    `json:"display_priority" validate:"gte=0"`
	Active          *bool  `json:"active,omitempty"`
}

type UpdateFactionRequest struct {
	Name            *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Acronym         *string `json:"acronym,omitempty" validate:"omitempty,min=2,max=10"`
	ColorHex        *string `json:"color_hex,omitempty" validate:"omitempty,hexcolor"`
	DisplayPriority *int    `json:"display_priority,omitempty" validate:"omitempty,gte=0"`
	Active          *bool   `json:"active,omitempty"`
}

func NewFactionHandler(factionRepo repository.FactionRepository) *FactionHandler {
	return &FactionHandler{factionRepo: factionRepo}
}

// GetFactions handles getting all factions
// @Summary List factions
// @Description Get list of factions with optional filtering
// @Tags Factions
// @Produce json
// @Param active query bool false "Filter by active status"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} utils.ErrorResponse
// @Router /factions [get]
func (h *FactionHandler) GetFactions(c *fiber.Ctx) error {
	// Parse query parameters
	activeParam := c.Query("active")
	var active *bool
	if activeParam != "" {
		activeBool, err := strconv.ParseBool(activeParam)
		if err != nil {
			return c.Status(400).JSON(utils.ErrorResponse{
				Error:   "Invalid active parameter",
				Code:    "VALIDATION_003",
				Message: "Active parameter must be a boolean value",
			})
		}
		active = &activeBool
	}

	factions, err := h.factionRepo.List(active)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse{
			Error:   "Failed to fetch factions",
			Code:    "FACTION_001",
			Message: "Unable to retrieve factions from database",
		})
	}

	return c.JSON(fiber.Map{
		"data": factions,
	})
}

// GetFaction handles getting a single faction by ID
// @Summary Get faction by ID
// @Description Get a specific faction by its ID
// @Tags Factions
// @Produce json
// @Param id path string true "Faction ID"
// @Success 200 {object} models.Faction
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /factions/{id} [get]
func (h *FactionHandler) GetFaction(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(utils.ErrorResponse{
			Error:   "Invalid faction ID",
			Code:    "VALIDATION_004",
			Message: "Faction ID must be a valid UUID",
		})
	}

	faction, err := h.factionRepo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse{
			Error:   "Faction not found",
			Code:    "FACTION_002",
			Message: "The requested faction was not found",
		})
	}

	return c.JSON(faction)
}

// CreateFaction handles creating a new faction
// @Summary Create faction
// @Description Create a new faction (Admin only)
// @Tags Factions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateFactionRequest true "Faction data"
// @Success 201 {object} models.Faction
// @Failure 400 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Router /factions [post]
func (h *FactionHandler) CreateFaction(c *fiber.Ctx) error {
	var req CreateFactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse{
			Error:   "Invalid request body",
			Code:    "VALIDATION_001",
			Message: "Request body must be valid JSON",
		})
	}

	// Validate input
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse{
			Error:   "Validation failed",
			Code:    "VALIDATION_002",
			Message: err.Error(),
		})
	}

	// Check if acronym already exists
	existingFaction, err := h.factionRepo.GetByAcronym(req.Acronym)
	if err == nil && existingFaction != nil {
		return c.Status(409).JSON(utils.ErrorResponse{
			Error:   "Acronym already exists",
			Code:    "FACTION_003",
			Message: "A faction with this acronym already exists",
		})
	}

	faction := &models.Faction{
		Name:            req.Name,
		Acronym:         req.Acronym,
		ColorHex:        req.ColorHex,
		DisplayPriority: req.DisplayPriority,
		Active:          true,
	}

	if req.Active != nil {
		faction.Active = *req.Active
	}

	if err := h.factionRepo.Create(faction); err != nil {
		return c.Status(500).JSON(utils.ErrorResponse{
			Error:   "Failed to create faction",
			Code:    "FACTION_004",
			Message: "Unable to create faction in database",
		})
	}

	return c.Status(201).JSON(faction)
}

// UpdateFaction handles updating an existing faction
// @Summary Update faction
// @Description Update an existing faction (Admin only)
// @Tags Factions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Faction ID"
// @Param request body UpdateFactionRequest true "Updated faction data"
// @Success 200 {object} models.Faction
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /factions/{id} [put]
func (h *FactionHandler) UpdateFaction(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(utils.ErrorResponse{
			Error:   "Invalid faction ID",
			Code:    "VALIDATION_004",
			Message: "Faction ID must be a valid UUID",
		})
	}

	// Get existing faction
	faction, err := h.factionRepo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse{
			Error:   "Faction not found",
			Code:    "FACTION_002",
			Message: "The requested faction was not found",
		})
	}

	var req UpdateFactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse{
			Error:   "Invalid request body",
			Code:    "VALIDATION_001",
			Message: "Request body must be valid JSON",
		})
	}

	// Validate input
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse{
			Error:   "Validation failed",
			Code:    "VALIDATION_002",
			Message: err.Error(),
		})
	}

	// Update fields if provided
	if req.Name != nil {
		faction.Name = *req.Name
	}
	if req.Acronym != nil {
		// Check if new acronym already exists (and it's not the same faction)
		existingFaction, err := h.factionRepo.GetByAcronym(*req.Acronym)
		if err == nil && existingFaction != nil && existingFaction.ID != faction.ID {
			return c.Status(409).JSON(utils.ErrorResponse{
				Error:   "Acronym already exists",
				Code:    "FACTION_003",
				Message: "A faction with this acronym already exists",
			})
		}
		faction.Acronym = *req.Acronym
	}
	if req.ColorHex != nil {
		faction.ColorHex = *req.ColorHex
	}
	if req.DisplayPriority != nil {
		faction.DisplayPriority = *req.DisplayPriority
	}
	if req.Active != nil {
		faction.Active = *req.Active
	}

	if err := h.factionRepo.Update(faction); err != nil {
		return c.Status(500).JSON(utils.ErrorResponse{
			Error:   "Failed to update faction",
			Code:    "FACTION_005",
			Message: "Unable to update faction in database",
		})
	}

	return c.JSON(faction)
}

// DeleteFaction handles deleting a faction
// @Summary Delete faction
// @Description Delete a faction (Admin only)
// @Tags Factions
// @Security BearerAuth
// @Param id path string true "Faction ID"
// @Success 204
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /factions/{id} [delete]
func (h *FactionHandler) DeleteFaction(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(utils.ErrorResponse{
			Error:   "Invalid faction ID",
			Code:    "VALIDATION_004",
			Message: "Faction ID must be a valid UUID",
		})
	}

	// Check if faction exists
	_, err = h.factionRepo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse{
			Error:   "Faction not found",
			Code:    "FACTION_002",
			Message: "The requested faction was not found",
		})
	}

	if err := h.factionRepo.Delete(id); err != nil {
		return c.Status(500).JSON(utils.ErrorResponse{
			Error:   "Failed to delete faction",
			Code:    "FACTION_006",
			Message: "Unable to delete faction from database",
		})
	}

	return c.SendStatus(204)
}