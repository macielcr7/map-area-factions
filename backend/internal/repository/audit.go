package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// GetLogs returns paginated audit logs with filters
func (r *AuditRepository) GetLogs(page, limit int, entity, action string, userID *uuid.UUID, fromDate, toDate *time.Time) ([]models.AuditLog, int, error) {
	var logs []models.AuditLog
	var total int64

	query := r.db.Model(&models.AuditLog{}).Preload("User")

	// Apply filters
	if entity != "" {
		query = query.Where("entity = ?", entity)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if fromDate != nil {
		query = query.Where("created_at >= ?", *fromDate)
	}
	if toDate != nil {
		query = query.Where("created_at <= ?", *toDate)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, int(total), nil
}

// GetByID returns a specific audit log by ID
func (r *AuditRepository) GetByID(id uuid.UUID) (*models.AuditLog, error) {
	var log models.AuditLog
	if err := r.db.Preload("User").First(&log, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// Create creates a new audit log entry
func (r *AuditRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

// GetStats returns audit statistics
func (r *AuditRepository) GetStats(days int) (map[string]interface{}, error) {
	fromDate := time.Now().AddDate(0, 0, -days)
	
	stats := make(map[string]interface{})
	
	// Total logs in period
	var totalLogs int64
	if err := r.db.Model(&models.AuditLog{}).Where("created_at >= ?", fromDate).Count(&totalLogs).Error; err != nil {
		return nil, err
	}
	stats["total_logs"] = totalLogs
	
	// Logs by action
	var actionStats []struct {
		Action string `json:"action"`
		Count  int    `json:"count"`
	}
	if err := r.db.Model(&models.AuditLog{}).
		Select("action, COUNT(*) as count").
		Where("created_at >= ?", fromDate).
		Group("action").
		Find(&actionStats).Error; err != nil {
		return nil, err
	}
	stats["by_action"] = actionStats
	
	// Logs by entity
	var entityStats []struct {
		Entity string `json:"entity"`
		Count  int    `json:"count"`
	}
	if err := r.db.Model(&models.AuditLog{}).
		Select("entity, COUNT(*) as count").
		Where("created_at >= ?", fromDate).
		Group("entity").
		Find(&entityStats).Error; err != nil {
		return nil, err
	}
	stats["by_entity"] = entityStats
	
	// Daily activity
	var dailyStats []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}
	if err := r.db.Model(&models.AuditLog{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", fromDate).
		Group("DATE(created_at)").
		Order("date").
		Find(&dailyStats).Error; err != nil {
		return nil, err
	}
	stats["daily_activity"] = dailyStats
	
	return stats, nil
}

// GetRecentActivities returns recent activities for dashboard
func (r *AuditRepository) GetRecentActivities(limit int) ([]RecentActivity, error) {
	var activities []RecentActivity
	
	query := `
		SELECT 
			al.id,
			al.user_id,
			u.name as user_name,
			al.entity,
			al.action,
			al.target,
			al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		ORDER BY al.created_at DESC
		LIMIT ?
	`
	
	if err := r.db.Raw(query, limit).Scan(&activities).Error; err != nil {
		return nil, err
	}
	
	return activities, nil
}

type RecentActivity struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	Entity    string    `json:"entity"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	CreatedAt time.Time `json:"created_at"`
}