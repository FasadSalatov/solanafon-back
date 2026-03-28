package models

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

// Referral — tracks referral relationships
type Referral struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	ReferrerID   uint      `gorm:"not null;index" json:"referrerId"`
	Referrer     User      `gorm:"foreignKey:ReferrerID" json:"-"`
	ReferredID   uint      `gorm:"not null;uniqueIndex" json:"referredId"`
	Referred     User      `gorm:"foreignKey:ReferredID" json:"referred,omitempty"`
	RegisteredAt time.Time `json:"registeredAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

// GenerateReferralCode creates a unique 8-char referral code
func GenerateReferralCode() string {
	b := make([]byte, 4)
	rand.Read(b)
	return strings.ToUpper(hex.EncodeToString(b))
}
