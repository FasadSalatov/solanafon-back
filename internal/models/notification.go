package models

import "time"

// Notification — user notification
type Notification struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Title     string    `gorm:"not null" json:"title"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	Type      string    `gorm:"not null" json:"type"` // app_moderation, security, system, transaction, promotion
	IsRead    bool      `gorm:"default:false" json:"isRead"`
	ActionURL string    `json:"actionUrl,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// PushToken — FCM/APNS push token registration
type PushToken struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	FCMToken  string    `gorm:"not null" json:"fcmToken"`
	DeviceID  string    `gorm:"not null" json:"deviceId"`
	Platform  string    `gorm:"default:android" json:"platform"` // android, ios
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
