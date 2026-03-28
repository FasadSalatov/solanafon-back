package models

import (
	"time"

	"gorm.io/gorm"
)

// NewsPost — news feed post by app developer
type NewsPost struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	AppID         uint           `gorm:"not null;index" json:"appId"`
	App           MiniApp        `gorm:"foreignKey:AppID" json:"app,omitempty"`
	Text          string         `gorm:"type:text;not null" json:"text"`
	ImageURL      string         `json:"imageUrl,omitempty"`
	CommentsCount int            `gorm:"default:0" json:"commentsCount"`
	LikesCount    int            `gorm:"default:0" json:"likesCount"`
	SharesCount   int            `gorm:"default:0" json:"sharesCount"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// NewsLike — post like
type NewsLike struct {
	PostID    uint      `gorm:"primaryKey;autoIncrement:false" json:"postId"`
	UserID    uint      `gorm:"primaryKey;autoIncrement:false" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

// NewsComment — post comment
type NewsComment struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	PostID    uint      `gorm:"not null;index" json:"postId"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}
