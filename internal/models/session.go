package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

// RefreshToken — JWT refresh token storage
type RefreshToken struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Token     string    `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expiresAt"`
	SessionID uint      `json:"sessionId"`
	CreatedAt time.Time `json:"createdAt"`
}

// Session — user device sessions
type Session struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"not null;index" json:"userId"`
	User       User           `gorm:"foreignKey:UserID" json:"-"`
	DeviceType string         `gorm:"default:unknown" json:"deviceType"` // android, ios, web
	DeviceName string         `json:"deviceName"`
	IPAddress  string         `json:"ipAddress"`
	Location   string         `json:"location"`
	LastActive time.Time      `json:"lastActive"`
	IsCurrent  bool           `gorm:"-" json:"isCurrent"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// GenerateRefreshToken creates a secure refresh token
func GenerateRefreshToken() string {
	b := make([]byte, 48)
	rand.Read(b)
	return hex.EncodeToString(b)
}
