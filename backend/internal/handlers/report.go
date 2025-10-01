package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"github.com/macielcr7/map-area-factions/backend/internal/repository"
	"github.com/macielcr7/map-area-factions/backend/internal/utils"
)

type ReportHandler struct {
	reportRepo *repository.ReportRepository
}

type CreateReportRequest struct {
	GeometryID  *uuid.UUID `json:"geometry_id" validate:"omitempty,uuid"`
	Type        string     `json:"type" validate:"required,oneof=inaccuracy inappropriate spam outdated other"`
	Description string     `json:"description" validate:"required,min=10,max=1000"`
	Attachments []string   `json:"attachments" validate:"omitempty,dive,url"`
	Location    *struct {
		Lat float64 `json:"lat" validate:"required,min=-90,max=90"`
		Lng float64 `json:"lng" validate:"required,min=-180,max=180"`
	} `json:"location" validate:"omitempty"`
}

type UpdateReportRequest struct {
	Status       string  `json:"status" validate:"omitempty,oneof=pending reviewing resolved rejected"`
	ReviewNotes  string  `json:"review_notes" validate:"omitempty,max=500"`
	AssignedToID *uuid.UUID `json:"assigned_to_id" validate:"omitempty,uuid"`
}

func NewReportHandler(reportRepo *repository.ReportRepository) *ReportHandler {
	return &ReportHandler{
		reportRepo: reportRepo,
	}
}

// GetReports returns paginated list of reports
func (h *ReportHandler) GetReports(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	status := c.Query("status", "")
	reportType := c.Query("type", "")
	assignedToStr := c.Query("assigned_to", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var assignedTo *uuid.UUID
	if assignedToStr != "" {
		if id, err := uuid.Parse(assignedToStr); err == nil {
			assignedTo = &id
		}
	}

	reports, total, err := h.reportRepo.GetReports(page, limit, status, reportType, assignedTo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "REPORT_001",
			"message": "Failed to fetch reports",
		})
	}

	return c.JSON(fiber.Map{
		"reports": reports,
		"pagination": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"pages":      (total + limit - 1) / limit,
			"has_next":   page*limit < total,
			"has_prev":   page > 1,
		},
		"filters": fiber.Map{
			"status":      status,
			"type":        reportType,
			"assigned_to": assignedToStr,
		},
	})
}

// GetReport returns a specific report by ID
func (h *ReportHandler) GetReport(c *fiber.Ctx) error {
	id := c.Params("id")
	reportID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "REPORT_002",
			"message": "Invalid report ID format",
		})
	}

	report, err := h.reportRepo.GetByID(reportID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "REPORT_003",
			"message": "Report not found",
		})
	}

	return c.JSON(fiber.Map{
		"report": report,
	})
}

// CreateReport creates a new report
func (h *ReportHandler) CreateReport(c *fiber.Ctx) error {
	var req CreateReportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "REPORT_004",
			"message": "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation Error",
			"code":    "REPORT_005",
			"message": "Invalid input data",
			"details": err,
		})
	}

	// Get user ID from context (could be anonymous)
	var userID *uuid.UUID
	if id, ok := c.Locals("user_id").(uuid.UUID); ok {
		userID = &id
	}

	report := &models.Report{
		UserID:      userID,
		GeometryID:  req.GeometryID,
		Type:        req.Type,
		Description: req.Description,
		Attachments: req.Attachments,
		Status:      "pending",
	}

	if req.Location != nil {
		report.ReporterLat = &req.Location.Lat
		report.ReporterLng = &req.Location.Lng
	}

	if err := h.reportRepo.Create(report); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "REPORT_006",
			"message": "Failed to create report",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Report created successfully",
		"report":  report,
	})
}

// UpdateReport updates an existing report (moderator/admin only)
func (h *ReportHandler) UpdateReport(c *fiber.Ctx) error {
	id := c.Params("id")
	reportID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "REPORT_007",
			"message": "Invalid report ID format",
		})
	}

	var req UpdateReportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "REPORT_008",
			"message": "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation Error",
			"code":    "REPORT_009",
			"message": "Invalid input data",
			"details": err,
		})
	}

	report, err := h.reportRepo.GetByID(reportID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "REPORT_010",
			"message": "Report not found",
		})
	}

	// Update fields if provided
	if req.Status != "" {
		report.Status = req.Status
	}
	if req.ReviewNotes != "" {
		report.ReviewNotes = &req.ReviewNotes
	}
	if req.AssignedToID != nil {
		report.AssignedToID = req.AssignedToID
	}

	// Set reviewer ID
	reviewerID := c.Locals("user_id").(uuid.UUID)
	report.ReviewedByID = &reviewerID

	if err := h.reportRepo.Update(report); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "REPORT_011",
			"message": "Failed to update report",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Report updated successfully",
		"report":  report,
	})
}

// DeleteReport deletes a report (admin only)
func (h *ReportHandler) DeleteReport(c *fiber.Ctx) error {
	id := c.Params("id")
	reportID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "REPORT_012",
			"message": "Invalid report ID format",
		})
	}

	report, err := h.reportRepo.GetByID(reportID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "REPORT_013",
			"message": "Report not found",
		})
	}

	if err := h.reportRepo.Delete(report); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "REPORT_014",
			"message": "Failed to delete report",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Report deleted successfully",
	})
}

// GetReportStats returns report statistics
func (h *ReportHandler) GetReportStats(c *fiber.Ctx) error {
	stats, err := h.reportRepo.GetStats()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "REPORT_015",
			"message": "Failed to fetch report statistics",
		})
	}

	return c.JSON(fiber.Map{
		"stats": stats,
	})
}