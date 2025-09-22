package repository

import (
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
	var geometries []models.Geometry
	var total int64

	query := r.db.Model(&models.Geometry{}).Preload("Faction")
	
	if factionID != "" {
		query = query.Where("faction_id = ?", factionID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&geometries).Error; err != nil {
		return nil, 0, err
	}

	return geometries, total, nil
}

func (r *GeometryRepository) GetByID(id uint) (*models.Geometry, error) {
	var geometry models.Geometry
	if err := r.db.Preload("Faction").First(&geometry, id).Error; err != nil {
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

func (r *GeometryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Geometry{}, id).Error
}

func (r *GeometryRepository) SearchNearby(lat, lng, radius float64, query string) ([]models.Geometry, error) {
	var geometries []models.Geometry
	
	// Using PostGIS ST_DWithin for spatial search
	// ST_DWithin(geography, ST_Point(lng, lat), radius_in_meters)
	sqlQuery := `
		SELECT * FROM geometries 
		WHERE ST_DWithin(
			geometry::geography, 
			ST_Point(?, ?)::geography, 
			?
		)
	`
	
	args := []interface{}{lng, lat, radius}
	
	if query != "" {
		sqlQuery += " AND (name ILIKE ? OR description ILIKE ?)"
		args = append(args, "%"+query+"%", "%"+query+"%")
	}
	
	if err := r.db.Raw(sqlQuery, args...).Scan(&geometries).Error; err != nil {
		return nil, err
	}
	
	return geometries, nil
}

func (r *GeometryRepository) GetByFactionID(factionID uint) ([]models.Geometry, error) {
	var geometries []models.Geometry
	if err := r.db.Where("faction_id = ? AND is_active = ?", factionID, true).Find(&geometries).Error; err != nil {
		return nil, err
	}
	return geometries, nil
}

func (r *GeometryRepository) GetInBounds(northEast, southWest map[string]float64) ([]models.Geometry, error) {
	var geometries []models.Geometry
	
	// Using PostGIS ST_Intersects with a bounding box
	sqlQuery := `
		SELECT * FROM geometries 
		WHERE ST_Intersects(
			geometry,
			ST_MakeEnvelope(?, ?, ?, ?, 4326)
		) AND is_active = true
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