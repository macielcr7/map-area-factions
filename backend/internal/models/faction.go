package models

type Faction struct {
	BaseModel
	Name            string `json:"name" gorm:"not null"`
	Acronym         string `json:"acronym" gorm:"uniqueIndex;not null"`
	ColorHex        string `json:"color_hex" gorm:"not null"`
	DisplayPriority int    `json:"display_priority" gorm:"default:0"`
	Active          bool   `json:"active" gorm:"default:true"`
}

// TableName specifies the table name for the Faction model
func (Faction) TableName() string {
	return "factions"
}