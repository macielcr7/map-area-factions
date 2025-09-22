package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/repository"
)

type AuditHandler struct {
	auditRepo *repository.AuditRepository
}

func NewAuditHandler(auditRepo *repository.AuditRepository) *AuditHandler {
	return &AuditHandler{
		auditRepo: auditRepo,
	}
}

// GetAuditLogs returns paginated list of audit logs
func (h *AuditHandler) GetAuditLogs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	entity := c.Query("entity", "")
	action := c.Query("action", "")
	userIDStr := c.Query("user_id", "")
	fromDateStr := c.Query("from_date", "")
	toDateStr := c.Query("to_date", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var userID *uuid.UUID
	if userIDStr != "" {
		if id, err := uuid.Parse(userIDStr); err == nil {
			userID = &id
		}
	}

	var fromDate, toDate *time.Time
	if fromDateStr != "" {
		if t, err := time.Parse(time.RFC3339, fromDateStr); err == nil {
			fromDate = &t
		}
	}
	if toDateStr != "" {
		if t, err := time.Parse(time.RFC3339, toDateStr); err == nil {
			toDate = &t
		}
	}

	logs, total, err := h.auditRepo.GetLogs(page, limit, entity, action, userID, fromDate, toDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "AUDIT_001",
			"message": "Failed to fetch audit logs",
		})
	}

	return c.JSON(fiber.Map{
		"logs": logs,
		"pagination": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"pages":      (total + limit - 1) / limit,
			"has_next":   page*limit < total,
			"has_prev":   page > 1,
		},
		"filters": fiber.Map{
			"entity":    entity,
			"action":    action,
			"user_id":   userIDStr,
			"from_date": fromDateStr,
			"to_date":   toDateStr,
		},
	})
}

// GetAuditLog returns a specific audit log by ID
func (h *AuditHandler) GetAuditLog(c *fiber.Ctx) error {
	id := c.Params("id")
	logID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "AUDIT_002",
			"message": "Invalid audit log ID format",
		})
	}

	log, err := h.auditRepo.GetByID(logID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "AUDIT_003",
			"message": "Audit log not found",
		})
	}

	return c.JSON(fiber.Map{
		"log": log,
	})
}

// GetAuditStats returns audit statistics
func (h *AuditHandler) GetAuditStats(c *fiber.Ctx) error {
	days, _ := strconv.Atoi(c.Query("days", "30"))
	if days < 1 || days > 365 {
		days = 30
	}

	stats, err := h.auditRepo.GetStats(days)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "AUDIT_004",
			"message": "Failed to fetch audit statistics",
		})
	}

	return c.JSON(fiber.Map{
		"stats": stats,
		"period_days": days,
	})
}