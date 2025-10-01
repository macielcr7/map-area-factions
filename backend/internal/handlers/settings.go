package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/macielcr7/map-area-factions/backend/internal/repository"
	"github.com/macielcr7/map-area-factions/backend/internal/utils"
)

type SettingsGeneral struct {
	SiteName        string `json:"site_name" validate:"required,min=2,max=150"`
	SiteDescription string `json:"site_description" validate:"required,min=4,max=255"`
	MaintenanceMode bool   `json:"maintenance_mode"`
	Timezone        string `json:"timezone" validate:"required"`
	DefaultLanguage string `json:"default_language" validate:"required"`
}

type SettingsMap struct {
	DefaultStyle     string  `json:"default_style" validate:"required"`
	DefaultZoom      float64 `json:"default_zoom" validate:"required,min=1,max=20"`
	DefaultCenterLat float64 `json:"default_center_lat" validate:"required,min=-90,max=90"`
	DefaultCenterLng float64 `json:"default_center_lng" validate:"required,min=-180,max=180"`
	MaxZoom          float64 `json:"max_zoom" validate:"required,min=1,max=22"`
	MinZoom          float64 `json:"min_zoom" validate:"required,min=0,max=22"`
}

type SettingsSecurity struct {
	SessionTimeoutMinutes    int  `json:"session_timeout_minutes" validate:"required,min=5,max=1440"`
	RequireEmailVerification bool `json:"require_email_verification"`
	EnableTwoFactor          bool `json:"enable_two_factor"`
	PasswordMinLength        int  `json:"password_min_length" validate:"required,min=6,max=64"`
	EnableAuditLogging       bool `json:"enable_audit_logging"`
}

type SettingsNotifications struct {
	EmailEnabled     bool   `json:"email_enabled"`
	EmailSMTPHost    string `json:"email_smtp_host" validate:"omitempty,hostname|ip"`
	EmailSMTPPort    int    `json:"email_smtp_port" validate:"omitempty,min=1,max=65535"`
	PushEnabled      bool   `json:"push_enabled"`
	WebhookURL       string `json:"webhook_url" validate:"omitempty,url"`
	NotifyNewReports bool   `json:"notify_new_reports"`
}

type SystemSettingsPayload struct {
	General       SettingsGeneral       `json:"general" validate:"required"`
	Map           SettingsMap           `json:"map" validate:"required"`
	Security      SettingsSecurity      `json:"security" validate:"required"`
	Notifications SettingsNotifications `json:"notifications" validate:"required"`
}

type SettingsHandler struct {
	repo            *repository.SettingsRepository
	defaultSettings SystemSettingsPayload
}

func NewSettingsHandler(repo *repository.SettingsRepository) *SettingsHandler {
	defaultPayload := SystemSettingsPayload{
		General: SettingsGeneral{
			SiteName:        "Map Area Factions",
			SiteDescription: "Sistema de mapeamento de áreas por facção",
			MaintenanceMode: false,
			Timezone:        "America/Sao_Paulo",
			DefaultLanguage: "pt-BR",
		},
		Map: SettingsMap{
			DefaultStyle:     "mapbox://styles/mapbox/streets-v11",
			DefaultZoom:      10,
			DefaultCenterLat: -3.7319,
			DefaultCenterLng: -38.5267,
			MaxZoom:          18,
			MinZoom:          8,
		},
		Security: SettingsSecurity{
			SessionTimeoutMinutes:    480,
			RequireEmailVerification: true,
			EnableTwoFactor:          false,
			PasswordMinLength:        8,
			EnableAuditLogging:       true,
		},
		Notifications: SettingsNotifications{
			EmailEnabled:     true,
			EmailSMTPHost:    "smtp.gmail.com",
			EmailSMTPPort:    587,
			PushEnabled:      false,
			WebhookURL:       "",
			NotifyNewReports: true,
		},
	}

	return &SettingsHandler{
		repo:            repo,
		defaultSettings: defaultPayload,
	}
}

func (h *SettingsHandler) GetSettings(c *fiber.Ctx) error {
	settings, err := h.repo.GetSettings(h.defaultSettings)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "SETTINGS_001",
			"message": "Failed to load settings",
		})
	}

	var payload SystemSettingsPayload
	if err := json.Unmarshal(settings.Data, &payload); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "SETTINGS_002",
			"message": "Failed to parse settings",
		})
	}

	return c.JSON(fiber.Map{
		"settings": payload,
	})
}

func (h *SettingsHandler) UpdateSettings(c *fiber.Ctx) error {
	var req SystemSettingsPayload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "SETTINGS_003",
			"message": "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation Error",
			"code":    "SETTINGS_004",
			"message": "Invalid input data",
			"details": err,
		})
	}

	settings, err := h.repo.GetSettings(h.defaultSettings)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "SETTINGS_005",
			"message": "Failed to load settings",
		})
	}

	if err := h.repo.UpdateSettings(settings, req); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "SETTINGS_006",
			"message": "Failed to update settings",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Settings updated successfully",
		"settings": req,
	})
}
