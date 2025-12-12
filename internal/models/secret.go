package models

import (
	"time"

	"gorm.io/gorm"
)

// SecretNumber - виртуальные номера для Secret Login
type SecretNumber struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Number      string         `gorm:"unique;not null" json:"number"` // Format: +999 (XXX) XXX-XX-XX
	IsPremium   bool           `gorm:"default:false" json:"isPremium"`
	PriceMP     int            `gorm:"not null" json:"priceMP"` // Price in Mana Points
	IsAvailable bool           `gorm:"default:true" json:"isAvailable"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// SecretAccess - активированный секретный доступ
type SecretAccess struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	UserID         uint           `gorm:"not null;uniqueIndex" json:"userId"`
	User           User           `gorm:"foreignKey:UserID" json:"-"`
	SecretNumberID uint           `gorm:"not null" json:"secretNumberId"`
	SecretNumber   SecretNumber   `gorm:"foreignKey:SecretNumberID" json:"secretNumber,omitempty"`
	ActivatedAt    time.Time      `gorm:"not null" json:"activatedAt"`
	IsActive       bool           `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
