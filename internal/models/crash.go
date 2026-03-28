package models

import "time"

// CrashReport — mobile crash log
type CrashReport struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    string    `json:"userId"`
	SessionID string    `json:"sessionId"`
	ErrorType string    `gorm:"not null" json:"errorType"` // crash, error, warning
	Message   string    `gorm:"type:text;not null" json:"message"`
	Stack     string    `gorm:"type:text" json:"stack"`
	UserAgent string    `json:"userAgent"`
	Source    string    `gorm:"default:android" json:"source"` // android, ios
	Metadata  string    `gorm:"type:jsonb" json:"metadata"`
	CreatedAt time.Time `json:"createdAt"`
}
