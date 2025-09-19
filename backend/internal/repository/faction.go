package repository

import (
	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"gorm.io/gorm"
)

type FactionRepository interface {
	Create(faction *models.Faction) error
	GetByID(id uuid.UUID) (*models.Faction, error)
	GetByAcronym(acronym string) (*models.Faction, error)
	Update(faction *models.Faction) error
	Delete(id uuid.UUID) error
	List(active *bool) ([]*models.Faction, error)
}

type factionRepository struct {
	db *gorm.DB
}

func NewFactionRepository(db *gorm.DB) FactionRepository {
	return &factionRepository{db: db}
}

func (r *factionRepository) Create(faction *models.Faction) error {
	return r.db.Create(faction).Error
}

func (r *factionRepository) GetByID(id uuid.UUID) (*models.Faction, error) {
	var faction models.Faction
	err := r.db.Where("id = ?", id).First(&faction).Error
	if err != nil {
		return nil, err
	}
	return &faction, nil
}

func (r *factionRepository) GetByAcronym(acronym string) (*models.Faction, error) {
	var faction models.Faction
	err := r.db.Where("acronym = ?", acronym).First(&faction).Error
	if err != nil {
		return nil, err
	}
	return &faction, nil
}

func (r *factionRepository) Update(faction *models.Faction) error {
	return r.db.Save(faction).Error
}

func (r *factionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Faction{}, id).Error
}

func (r *factionRepository) List(active *bool) ([]*models.Faction, error) {
	var factions []*models.Faction
	query := r.db.Model(&models.Faction{})

	if active != nil {
		query = query.Where("active = ?", *active)
	}

	err := query.Order("display_priority DESC, name ASC").Find(&factions).Error
	return factions, err
}