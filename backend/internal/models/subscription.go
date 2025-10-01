package models

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusPending   SubscriptionStatus = "pending"
)

type Subscription struct {
	BaseModel
	UserID        *uuid.UUID         `json:"user_id" gorm:"type:uuid"`
	Plan          string             `json:"plan" gorm:"not null"`
	Status        SubscriptionStatus `json:"status" gorm:"default:'pending'"`
	StartDate     time.Time          `json:"start_date" gorm:"not null"`
	EndDate       time.Time          `json:"end_date" gorm:"not null"`
	Price         float64            `json:"price" gorm:"not null"`
	Currency      string             `json:"currency" gorm:"default:'BRL'"`
	PaymentMethod string             `json:"payment_method"`
	PaymentID     string             `json:"payment_id"` // External payment system ID
	AutoRenew     bool               `json:"auto_renew" gorm:"default:false"`
	CancelledAt   *time.Time         `json:"cancelled_at"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for the Subscription model
func (Subscription) TableName() string {
	return "subscriptions"
}
