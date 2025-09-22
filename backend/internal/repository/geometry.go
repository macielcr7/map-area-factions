package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/macielcr7/map-area-factions/backend/internal/models"
)

type GeometryRepository struct {
	db *gorm.DB
}

func NewGeometryRepository(db *gorm.DB) *GeometryRepository {
	return &GeometryRepository{db: db}
}

func (r *GeometryRepository) GetAll(page, limit int, factionID string) ([]models.Geometry, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	var (
		geometries []models.Geometry
		total      int64
	)

	query := r.db.Model(&models.Geometry{}).
		Preload("Faction").
		Preload("Region").
		Preload("Author")

	if factionID != "" {
		query = query.Where("faction_id = ?", factionID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&geometries).Error; err != nil {
		return nil, 0, err
	}

	return geometries, total, nil
}

func (r *GeometryRepository) GetByID(id uuid.UUID) (*models.Geometry, error) {
	var geometry models.Geometry
	if err := r.db.Preload("Faction").Preload("Region").Preload("Author").First(&geometry, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &geometry, nil
}

func (r *GeometryRepository) Create(geometry *models.Geometry) error {
	return r.db.Create(geometry).Error
}

func (r *GeometryRepository) Update(geometry *models.Geometry) error {
	return r.db.Save(geometry).Error
}

func (r *GeometryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Geometry{}, "id = ?", id).Error
}

func (r *GeometryRepository) SearchNearby(lat, lng, radius float64, query string) ([]models.Geometry, error) {
	var geometries []models.Geometry

	sqlQuery := `
        SELECT * FROM geometries
        WHERE ST_DWithin(
            geom::geography,
            ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
            ?
        )
    `

	args := []interface{}{lng, lat, radius}

	if query != "" {
		sqlQuery += " AND (notes ILIKE ? OR source ILIKE ?)"
		like := "%" + query + "%"
		args = append(args, like, like)
	}

	if err := r.db.Raw(sqlQuery, args...).Scan(&geometries).Error; err != nil {
		return nil, err
	}

	return geometries, nil
}

func (r *GeometryRepository) GetByFactionID(factionID uuid.UUID) ([]models.Geometry, error) {
	var geometries []models.Geometry
	if err := r.db.Where("faction_id = ?", factionID).Find(&geometries).Error; err != nil {
		return nil, err
	}
	return geometries, nil
}

func (r *GeometryRepository) GetInBounds(northEast, southWest map[string]float64) ([]models.Geometry, error) {
	var geometries []models.Geometry

	sqlQuery := `
        SELECT * FROM geometries
        WHERE ST_Intersects(
            geom,
            ST_MakeEnvelope(?, ?, ?, ?, 4326)
        )
    `

	if err := r.db.Raw(
		sqlQuery,
		southWest["lng"], southWest["lat"],
		northEast["lng"], northEast["lat"],
	).Scan(&geometries).Error; err != nil {
		return nil, err
	}

	return geometries, nil
}
