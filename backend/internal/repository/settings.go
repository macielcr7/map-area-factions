package repository

import (
	"encoding/json"
	"errors"

	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) GetSettings(defaultPayload interface{}) (*models.SystemSettings, error) {
	var settings models.SystemSettings
	err := r.db.First(&settings).Error
	if err == nil {
		return &settings, nil
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	bytes, marshalErr := json.Marshal(defaultPayload)
	if marshalErr != nil {
		return nil, marshalErr
	}

	settings = models.SystemSettings{
		Data: datatypes.JSON(bytes),
	}

	if createErr := r.db.Create(&settings).Error; createErr != nil {
		return nil, createErr
	}

	return &settings, nil
}

func (r *SettingsRepository) UpdateSettings(settings *models.SystemSettings, payload interface{}) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	settings.Data = datatypes.JSON(bytes)
	return r.db.Save(settings).Error
}
