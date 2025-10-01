package repository

import (
	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"gorm.io/gorm"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetReports returns paginated reports with filters
func (r *ReportRepository) GetReports(page, limit int, status, reportType string, assignedTo *uuid.UUID) ([]models.Report, int, error) {
	var reports []models.Report
	var total int64

	query := r.db.Model(&models.Report{}).
		Preload("User").
		Preload("Geometry").
		Preload("AssignedTo").
		Preload("ReviewedBy")

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if reportType != "" {
		query = query.Where("type = ?", reportType)
	}
	if assignedTo != nil {
		query = query.Where("assigned_to_id = ?", *assignedTo)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, int(total), nil
}

// GetByID returns a specific report by ID
func (r *ReportRepository) GetByID(id uuid.UUID) (*models.Report, error) {
	var report models.Report
	if err := r.db.Preload("User").
		Preload("Geometry").
		Preload("AssignedTo").
		Preload("ReviewedBy").
		First(&report, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

// Create creates a new report
func (r *ReportRepository) Create(report *models.Report) error {
	return r.db.Create(report).Error
}

// Update updates an existing report
func (r *ReportRepository) Update(report *models.Report) error {
	return r.db.Save(report).Error
}

// Delete soft deletes a report
func (r *ReportRepository) Delete(report *models.Report) error {
	return r.db.Delete(report).Error
}

// GetStats returns report statistics
func (r *ReportRepository) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total reports
	var totalReports int64
	if err := r.db.Model(&models.Report{}).Count(&totalReports).Error; err != nil {
		return nil, err
	}
	stats["total_reports"] = totalReports
	
	// Reports by status
	var statusStats []struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}
	if err := r.db.Model(&models.Report{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&statusStats).Error; err != nil {
		return nil, err
	}
	stats["by_status"] = statusStats
	
	// Reports by type
	var typeStats []struct {
		Type  string `json:"type"`
		Count int    `json:"count"`
	}
	if err := r.db.Model(&models.Report{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Find(&typeStats).Error; err != nil {
		return nil, err
	}
	stats["by_type"] = typeStats
	
	// Pending reports (need attention)
	var pendingReports int64
	if err := r.db.Model(&models.Report{}).
		Where("status = ?", "pending").
		Count(&pendingReports).Error; err != nil {
		return nil, err
	}
	stats["pending_reports"] = pendingReports
	
	return stats, nil
}

func (r *ReportRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Report{}).Count(&count).Error
	return count, err
}