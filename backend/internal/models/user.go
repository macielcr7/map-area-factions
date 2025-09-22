package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin        UserRole = "admin"
	RoleModerator    UserRole = "moderator"
	RoleCollaborator UserRole = "collaborator"
	RoleCitizen      UserRole = "citizen"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      string         `json:"role" gorm:"not null;default:'citizen'"`
	Active    bool           `json:"active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

type AuditLog struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     *uuid.UUID `json:"user_id" gorm:"type:uuid"`
	User       *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Entity     string     `json:"entity" gorm:"not null"`
	EntityID   *uuid.UUID `json:"entity_id" gorm:"type:uuid"`
	Action     string     `json:"action" gorm:"not null"`
	Changes    string     `json:"changes" gorm:"type:text"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	CreatedAt  time.Time  `json:"created_at"`
}

type Report struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        *uuid.UUID `json:"user_id" gorm:"type:uuid"`
	User          *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	GeometryID    *uuid.UUID `json:"geometry_id" gorm:"type:uuid"`
	Geometry      *Geometry  `json:"geometry,omitempty" gorm:"foreignKey:GeometryID"`
	Type          string     `json:"type" gorm:"not null"`
	Description   string     `json:"description" gorm:"not null"`
	Attachments   []string   `json:"attachments" gorm:"type:text[]"`
	Status        string     `json:"status" gorm:"not null;default:'pending'"`
	ReviewNotes   *string    `json:"review_notes"`
	AssignedToID  *uuid.UUID `json:"assigned_to_id" gorm:"type:uuid"`
	AssignedTo    *User      `json:"assigned_to,omitempty" gorm:"foreignKey:AssignedToID"`
	ReviewedByID  *uuid.UUID `json:"reviewed_by_id" gorm:"type:uuid"`
	ReviewedBy    *User      `json:"reviewed_by,omitempty" gorm:"foreignKey:ReviewedByID"`
	ReporterLat   *float64   `json:"reporter_lat"`
	ReporterLng   *float64   `json:"reporter_lng"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}