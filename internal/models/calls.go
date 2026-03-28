package models

import (
	"crypto/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CallRoom — WebRTC call room
type CallRoom struct {
	ID           uint              `gorm:"primarykey" json:"id"`
	RoomCode     string            `gorm:"uniqueIndex;not null" json:"roomCode"` // ABC-DEF format
	Type         string            `gorm:"default:direct" json:"type"`           // direct, conference
	Status       string            `gorm:"default:waiting" json:"status"`        // waiting, active, ended
	CreatedBy    uint              `gorm:"not null" json:"createdBy"`
	Creator      User              `gorm:"foreignKey:CreatedBy" json:"-"`
	StartedAt    *time.Time        `json:"startedAt"`
	EndedAt      *time.Time        `json:"endedAt"`
	Duration     int               `gorm:"default:0" json:"duration"` // seconds
	Participants []CallParticipant `gorm:"foreignKey:RoomID" json:"participants,omitempty"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt    `gorm:"index" json:"-"`
}

// CallParticipant — participant in a call room
type CallParticipant struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	RoomID    uint       `gorm:"not null;index" json:"roomId"`
	UserID    uint       `gorm:"not null" json:"userId"`
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status    string     `gorm:"default:invited" json:"status"` // invited, joined, left
	JoinedAt  *time.Time `json:"joinedAt"`
	IsMuted   bool       `gorm:"default:false" json:"isMuted"`
	IsVideoOn bool       `gorm:"default:true" json:"isVideoOn"`
	IsAudioOn bool       `gorm:"default:true" json:"isAudioOn"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// GenerateRoomCode creates a unique room code like "ABC-DEF"
func GenerateRoomCode() string {
	b := make([]byte, 3)
	rand.Read(b)
	part1 := fmt.Sprintf("%c%c%c", 'A'+b[0]%26, 'A'+b[1]%26, 'A'+b[2]%26)
	rand.Read(b)
	part2 := fmt.Sprintf("%c%c%c", 'A'+b[0]%26, 'A'+b[1]%26, 'A'+b[2]%26)
	return part1 + "-" + part2
}
