package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"unique;not null" json:"name"`
	Slug        string         `gorm:"unique;not null" json:"slug"`
	Description string         `json:"description"`
	Icon        string         `json:"icon"`
	Order       int            `gorm:"default:0" json:"order"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// ModerationStatus - статус модерации приложения
type ModerationStatus string

const (
	ModerationPending  ModerationStatus = "pending"
	ModerationApproved ModerationStatus = "approved"
	ModerationRejected ModerationStatus = "rejected"
)

type MiniApp struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Subtitle    string         `json:"subtitle"`
	Icon        string         `json:"icon"` // Emoji
	CategoryID  uint           `gorm:"not null" json:"categoryId"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	IsSecret    bool           `gorm:"default:false" json:"isSecret"`
	UsersCount  int            `gorm:"default:0" json:"usersCount"`
	IsVerified  bool           `gorm:"default:false" json:"isVerified"`
	Description string         `gorm:"type:text" json:"description"`

	// App URL (optional - if empty, works as bot-only)
	URL string `json:"url,omitempty"`

	// Creator/Developer info
	CreatorID uint  `gorm:"index" json:"creatorId"`
	Creator   *User `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`

	// Developer API credentials (for bot functionality)
	APIToken   string `gorm:"unique" json:"-"`                    // Secret token for API calls
	WebhookURL string `json:"webhookUrl,omitempty"`               // URL to receive messages
	BotUsername string `gorm:"unique" json:"botUsername,omitempty"` // @myappbot style username

	// Bot welcome message (shown on /start)
	WelcomeMessage string `gorm:"type:text" json:"welcomeMessage,omitempty"`

	// Moderation
	ModerationStatus ModerationStatus `gorm:"default:pending" json:"moderationStatus"`
	ModerationNote   string           `json:"moderationNote,omitempty"`
	ModeratedAt      *time.Time       `json:"moderatedAt,omitempty"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// GenerateAPIToken creates a secure API token for the app
func GenerateAPIToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// AppUser - tracks which users use which apps (conversations)
type AppUser struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_app" json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	AppID     uint      `gorm:"not null;uniqueIndex:idx_user_app" json:"appId"`
	App       MiniApp   `gorm:"foreignKey:AppID" json:"app,omitempty"`
	LastUsed  time.Time `json:"lastUsed"`
	CreatedAt time.Time `json:"createdAt"`
}

// AppMessage - messages in app chat
type AppMessage struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	AppID     uint      `gorm:"not null;index" json:"appId"`
	App       MiniApp   `gorm:"foreignKey:AppID" json:"-"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IsFromBot bool      `gorm:"default:false" json:"isFromBot"`
	IsRead    bool      `gorm:"default:false" json:"isRead"`

	// Message metadata
	MessageType string `gorm:"default:text" json:"messageType"` // text, image, button, etc.
	Metadata    string `gorm:"type:jsonb" json:"metadata,omitempty"` // JSON for buttons, images, etc.

	CreatedAt time.Time `json:"createdAt"`
}

// BotCommand - predefined commands for the bot
type BotCommand struct {
	ID          uint    `gorm:"primarykey" json:"id"`
	AppID       uint    `gorm:"not null;index" json:"appId"`
	App         MiniApp `gorm:"foreignKey:AppID" json:"-"`
	Command     string  `gorm:"not null" json:"command"`     // e.g., "/start", "/help"
	Description string  `json:"description"`                  // Command description
	Response    string  `gorm:"type:text" json:"response"`   // Auto-response text
	IsEnabled   bool    `gorm:"default:true" json:"isEnabled"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// WebhookLog - logs of webhook calls
type WebhookLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	AppID      uint      `gorm:"not null;index" json:"appId"`
	URL        string    `json:"url"`
	Method     string    `json:"method"`
	StatusCode int       `json:"statusCode"`
	Request    string    `gorm:"type:text" json:"request"`
	Response   string    `gorm:"type:text" json:"response"`
	Duration   int       `json:"duration"` // milliseconds
	CreatedAt  time.Time `json:"createdAt"`
}

// ConversationState - tracks user's conversation state with BotFather
type ConversationState struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex" json:"userId"`
	State     string    `gorm:"not null" json:"state"`
	Data      string    `gorm:"type:jsonb" json:"data"`
	UpdatedAt time.Time `json:"updatedAt"`
}
