package repository

import (
	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"gorm.io/gorm"
)

type IncidentRepository struct {
	db *gorm.DB
}

func NewIncidentRepository(db *gorm.DB) *IncidentRepository {
	return &IncidentRepository{db: db}
}

func (r *IncidentRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Incident{}).Count(&count).Error
	return count, err
}

func (r *IncidentRepository) Create(incident *models.Incident) error {
	return r.db.Create(incident).Error
}

func (r *IncidentRepository) GetByID(id uuid.UUID) (*models.Incident, error) {
	var incident models.Incident
	err := r.db.Where("id = ?", id).First(&incident).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

func (r *IncidentRepository) Update(incident *models.Incident) error {
	return r.db.Save(incident).Error
}

func (r *IncidentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Incident{}, id).Error
}

func (r *IncidentRepository) List(limit, offset int, status string) ([]*models.Incident, int64, error) {
	var incidents []*models.Incident
	var count int64

	query := r.db.Model(&models.Incident{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&incidents).Error; err != nil {
		return nil, 0, err
	}

	return incidents, count, nil
}
