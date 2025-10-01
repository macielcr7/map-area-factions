package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type SystemSettings struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Data      datatypes.JSON `json:"data" gorm:"type:jsonb;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (SystemSettings) TableName() string {
	return "system_settings"
}
