package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Email           string         `gorm:"unique;not null" json:"email"`
	Name            string         `gorm:"default:'Solana User'" json:"name"`
	Avatar          string         `json:"avatar"`
	ManaPoints      int            `gorm:"default:0" json:"manaPoints"`
	HasSecretAccess bool           `gorm:"default:false" json:"hasSecretAccess"`
	Language        string         `gorm:"default:en" json:"language"`

	// Settings
	NotificationsEnabled bool `gorm:"default:true" json:"notificationsEnabled"`
	TwoFactorEnabled     bool `gorm:"default:false" json:"twoFactorEnabled"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserStats - computed stats for profile
type UserStats struct {
	AppsCount         int `json:"appsCount"`
	TransactionsCount int `json:"transactionsCount"`
	NFTsCount         int `json:"nftsCount"`
}

type OTP struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Email     string    `gorm:"not null;index" json:"email"`
	Code      string    `gorm:"not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expiresAt"`
	Verified  bool      `gorm:"default:false" json:"verified"`
	CreatedAt time.Time `json:"createdAt"`
}

// ManaTransaction - history of MP transactions
type ManaTransaction struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"userId"`
	User        User      `gorm:"foreignKey:UserID" json:"-"`
	Amount      int       `gorm:"not null" json:"amount"` // positive = credit, negative = debit
	Type        string    `gorm:"not null" json:"type"`   // topup, purchase, reward, refund
	Description string    `json:"description"`
	ReferenceID *uint     `json:"referenceId,omitempty"` // Related entity ID (app, secret number, etc.)
	CreatedAt   time.Time `json:"createdAt"`
}
