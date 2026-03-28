package models

import "time"

// FAQ — frequently asked questions
type FAQ struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Question string `gorm:"type:text;not null" json:"question"`
	Answer   string `gorm:"type:text;not null" json:"answer"`
	Category string `gorm:"not null" json:"category"` // wallet, general, developer, account
	Language string `gorm:"default:en" json:"language"`
	Order    int    `gorm:"default:0" json:"order"`
}

// SupportTicket — user support ticket
type SupportTicket struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Subject   string    `gorm:"not null" json:"subject"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Category  string    `gorm:"default:general" json:"category"` // wallet, general, developer, account
	Status    string    `gorm:"default:open" json:"status"`      // open, in_progress, resolved, closed
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// LegalDocument — terms of service, privacy policy
type LegalDocument struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Type        string    `gorm:"uniqueIndex;not null" json:"type"` // terms, privacy
	Content     string    `gorm:"type:text;not null" json:"content"`
	Version     string    `gorm:"default:1.0" json:"version"`
	LastUpdated string    `json:"lastUpdated"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
