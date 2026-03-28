package models

import (
	"time"

	"gorm.io/gorm"
)

// Conversation — chat conversation between user and app
type Conversation struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	AppID        uint           `gorm:"not null;index" json:"appId"`
	App          MiniApp        `gorm:"foreignKey:AppID" json:"app,omitempty"`
	UserID       uint           `gorm:"not null;index" json:"userId"`
	User         User           `gorm:"foreignKey:UserID" json:"-"`
	UnreadCount  int            `gorm:"default:0" json:"unreadCount"`
	IsActive     bool           `gorm:"default:true" json:"isActive"`
	LastMessageAt *time.Time    `json:"lastMessageAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// ChatMessage — message in a conversation
type ChatMessage struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	ConversationID uint      `gorm:"not null;index" json:"conversationId"`
	AppID          uint      `gorm:"not null;index" json:"appId"`
	SenderID       string    `gorm:"not null" json:"senderId"`   // user ID or "bot"
	SenderType     string    `gorm:"not null" json:"senderType"` // "user", "bot", "system"
	Content        string    `gorm:"type:jsonb;not null" json:"content"`
	Status         string    `gorm:"default:sent" json:"status"` // sending, sent, delivered, read, failed
	ReplyToID      *uint     `json:"replyToId"`
	Metadata       string    `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}
