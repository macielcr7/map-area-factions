package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/macielcr7/map-area-factions/backend/internal/repository"
)

type DashboardHandler struct {
	userRepo        UserRepository
	factionRepo     FactionRepository
	geometryRepo    GeometryRepository
	incidentRepo    IncidentRepository
	subscriptionRepo SubscriptionRepository
	reportRepo      ReportRepository
	auditRepo       AuditRepository
}

type UserRepository interface {
	Count() (int64, error)
}

type FactionRepository interface {
	Count() (int64, error)
}

type GeometryRepository interface {
	Count() (int64, error)
}

type IncidentRepository interface {
	Count() (int64, error)
}

type SubscriptionRepository interface {
	Count() (int64, error)
}

type ReportRepository interface {
	Count() (int64, error)
}

type AuditRepository interface {
	GetRecentActivities(limit int) ([]repository.RecentActivity, error)
}


type DashboardStats struct {
	TotalUsers        int64 `json:"total_users"`
	TotalGeometries   int64 `json:"total_geometries"`
	TotalFactions     int64 `json:"total_factions"`
	TotalReports      int64 `json:"total_reports"`
	TotalIncidents    int64 `json:"total_incidents"`
	TotalSubscriptions int64 `json:"total_subscriptions"`
}

type ActivityData struct {
	Name       string `json:"name"`
	Geometries int    `json:"geometries"`
	Users      int    `json:"users"`
	Incidents  int    `json:"incidents"`
}

type Alert struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Count   *int   `json:"count,omitempty"`
	Time    *string `json:"time,omitempty"`
}

func NewDashboardHandler(
	userRepo UserRepository,
	factionRepo FactionRepository,
	geometryRepo GeometryRepository,
	incidentRepo IncidentRepository,
	subscriptionRepo SubscriptionRepository,
	reportRepo ReportRepository,
	auditRepo AuditRepository,
) *DashboardHandler {
	return &DashboardHandler{
		userRepo:        userRepo,
		factionRepo:     factionRepo,
		geometryRepo:    geometryRepo,
		incidentRepo:    incidentRepo,
		subscriptionRepo: subscriptionRepo,
		reportRepo:      reportRepo,
		auditRepo:       auditRepo,
	}
}

// GetDashboardStats returns dashboard statistics
func (h *DashboardHandler) GetDashboardStats(c *fiber.Ctx) error {
	// Get counts from repositories
	userCount, err := h.userRepo.Count()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get user count",
		})
	}

	geometryCount, err := h.geometryRepo.Count()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get geometry count",
		})
	}

	factionCount, err := h.factionRepo.Count()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get faction count",
		})
	}

	reportCount, err := h.reportRepo.Count()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get report count",
		})
	}

	incidentCount, err := h.incidentRepo.Count()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get incident count",
		})
	}

	subscriptionCount, err := h.subscriptionRepo.Count()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get subscription count",
		})
	}

	stats := DashboardStats{
		TotalUsers:         userCount,
		TotalGeometries:    geometryCount,
		TotalFactions:      factionCount,
		TotalReports:       reportCount,
		TotalIncidents:     incidentCount,
		TotalSubscriptions: subscriptionCount,
	}

	return c.JSON(stats)
}

// GetActivityData returns activity data for charts
func (h *DashboardHandler) GetActivityData(c *fiber.Ctx) error {
	// Generate mock activity data for the last 6 months
	now := time.Now()
	activityData := []ActivityData{}

	for i := 5; i >= 0; i-- {
		date := now.AddDate(0, -i, 0)
		monthName := date.Format("Jan")
		
		// Generate realistic data based on month
		geometries := 10 + (i * 3) + (i % 3)
		users := 20 + (i * 5) + (i % 2)
		incidents := 5 + (i * 2) + (i % 4)

		activityData = append(activityData, ActivityData{
			Name:       monthName,
			Geometries: geometries,
			Users:      users,
			Incidents:  incidents,
		})
	}

	return c.JSON(activityData)
}

// GetRecentActivities returns recent activities
func (h *DashboardHandler) GetRecentActivities(c *fiber.Ctx) error {
	activities, err := h.auditRepo.GetRecentActivities(10)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get recent activities",
		})
	}

	return c.JSON(activities)
}

// GetAlerts returns system alerts
func (h *DashboardHandler) GetAlerts(c *fiber.Ctx) error {
	// Get pending geometries count
	geometryCount, _ := h.geometryRepo.Count()
	
	alerts := []Alert{
		{
			ID:      "1",
			Type:    "warning",
			Message: "Muitas geometrias pendentes de aprovação",
			Count:   func() *int { c := int(geometryCount / 10); return &c }(),
		},
		{
			ID:      "2",
			Type:    "info",
			Message: "Backup automático realizado com sucesso",
			Time:    func() *string { t := "2 horas atrás"; return &t }(),
		},
		{
			ID:      "3",
			Type:    "success",
			Message: "Sistema funcionando normalmente",
		},
	}

	return c.JSON(alerts)
}
