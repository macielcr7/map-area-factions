package models

import (
	"time"

	"github.com/google/uuid"
)

type IncidentStatus string

const (
	IncidentStatusActive   IncidentStatus = "active"
	IncidentStatusResolved IncidentStatus = "resolved"
	IncidentStatusExpired  IncidentStatus = "expired"
)

type Incident struct {
	BaseModel
	Type        string         `json:"type" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	Location    string         `json:"location" gorm:"not null"`
	Latitude    float64        `json:"latitude" gorm:"not null"`
	Longitude   float64        `json:"longitude" gorm:"not null"`
	Geom        interface{}    `json:"-" gorm:"type:geometry(POINT,4326);not null"`
	Status      IncidentStatus `json:"status" gorm:"default:'active'"`
	Severity    int            `json:"severity" gorm:"default:1"` // 1-5 scale
	Source      string         `json:"source"`
	AuthorID    *uuid.UUID     `json:"author_id" gorm:"type:uuid"`
	ResolvedAt  *time.Time     `json:"resolved_at"`
	ExpiresAt   *time.Time     `json:"expires_at"`

	// Relations
	Author *User `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
}

// TableName specifies the table name for the Incident model
func (Incident) TableName() string {
	return "incidents"
}
