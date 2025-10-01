package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type GeometryStatus string
type RiskLevel string

const (
	GeometryStatusDraft     GeometryStatus = "draft"
	GeometryStatusReview    GeometryStatus = "review"
	GeometryStatusPublished GeometryStatus = "published"

	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
)

type Region struct {
	BaseModel
	Name         string         `json:"name" gorm:"not null"`
	Type         string         `json:"type" gorm:"not null"` // neighborhood, zone, custom
	Municipality string         `json:"municipality" gorm:"not null"`
	State        string         `json:"state" gorm:"not null"`
	Centroid     interface{}    `json:"centroid" gorm:"type:geometry(POINT,4326)"`
	Bounds       interface{}    `json:"bounds" gorm:"type:geometry(POLYGON,4326)"`
	Status       GeometryStatus `json:"status" gorm:"default:'draft'"`
}

// TableName specifies the table name for the Region model
func (Region) TableName() string {
	return "regions"
}

type GeoJSON map[string]interface{}

// Value implements the driver.Valuer interface for GeoJSON
func (g GeoJSON) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// Scan implements the sql.Scanner interface for GeoJSON
func (g *GeoJSON) Scan(value interface{}) error {
	if value == nil {
		*g = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, g)
}

type Geometry struct {
	BaseModel
	RegionID       *uuid.UUID     `json:"region_id" gorm:"type:uuid"`
	GeoJSON        GeoJSON        `json:"geojson" gorm:"type:jsonb;not null"`
	GeometryType   string         `json:"geometry_type" gorm:"not null"` // polygon, polyline
	RiskLevel      RiskLevel      `json:"risk_level" gorm:"default:'low'"`
	FactionID      *uuid.UUID     `json:"faction_id" gorm:"type:uuid"`
	ValidityStart  *time.Time     `json:"validity_start"`
	ValidityEnd    *time.Time     `json:"validity_end"`
	Source         string         `json:"source"`
	Notes          string         `json:"notes"`
	Status         GeometryStatus `json:"status" gorm:"default:'draft'"`
	Geom           interface{}    `json:"-" gorm:"type:geometry;not null"`
	AuthorID       *uuid.UUID     `json:"author_id" gorm:"type:uuid"`

	// Relations
	Region   *Region  `json:"region,omitempty" gorm:"foreignKey:RegionID"`
	Faction  *Faction `json:"faction,omitempty" gorm:"foreignKey:FactionID"`
	Author   *User    `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
}

// TableName specifies the table name for the Geometry model
func (Geometry) TableName() string {
	return "geometries"
}