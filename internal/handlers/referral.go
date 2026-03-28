package handlers

import (
	"fmt"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ReferralHandler handles /api/referral/* endpoints
type ReferralHandler struct {
	db *gorm.DB
}

func NewReferralHandler(db *gorm.DB) *ReferralHandler {
	return &ReferralHandler{db: db}
}

// GetReferralInfo — GET /api/referral
func (h *ReferralHandler) GetReferralInfo(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var user models.User
	h.db.First(&user, userID)

	// Ensure user has a referral code
	if user.ReferralCode == "" {
		user.ReferralCode = models.GenerateReferralCode()
		h.db.Save(&user)
	}

	var referrals []models.Referral
	h.db.Where("referrer_id = ?", userID).Preload("Referred").Find(&referrals)

	refs := make([]fiber.Map, len(referrals))
	for i, r := range referrals {
		refs[i] = fiber.Map{
			"id": fmt.Sprintf("user_%d", r.ReferredID),
			"displayName":  r.Referred.GetDisplayName(),
			"avatarUrl":    r.Referred.GetAvatarURL(),
			"registeredAt": r.RegisteredAt,
		}
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"code":       user.ReferralCode,
			"link":       fmt.Sprintf("https://solafon.com/ref/%s", user.ReferralCode),
			"totalCount": len(referrals),
			"referrals":  refs,
		},
	})
}

// ValidateCode — GET /api/referral/validate/:code
func (h *ReferralHandler) ValidateCode(c *fiber.Ctx) error {
	code := c.Params("code")
	var user models.User
	if err := h.db.Where("referral_code = ?", code).First(&user).Error; err != nil {
		return c.JSON(fiber.Map{"data": fiber.Map{"valid": false}})
	}
	return c.JSON(fiber.Map{"data": fiber.Map{"valid": true, "referrerName": user.GetDisplayName()}})
}
