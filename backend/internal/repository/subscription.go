package repository

import (
	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Subscription{}).Count(&count).Error
	return count, err
}

func (r *SubscriptionRepository) Create(subscription *models.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Where("id = ?", id).First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepository) GetByUserID(userID uuid.UUID) ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.Where("user_id = ?", userID).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *SubscriptionRepository) Update(subscription *models.Subscription) error {
	return r.db.Save(subscription).Error
}

func (r *SubscriptionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Subscription{}, id).Error
}

func (r *SubscriptionRepository) List(limit, offset int, status string) ([]*models.Subscription, int64, error) {
	var subscriptions []*models.Subscription
	var count int64

	query := r.db.Model(&models.Subscription{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&subscriptions).Error; err != nil {
		return nil, 0, err
	}

	return subscriptions, count, nil
}
