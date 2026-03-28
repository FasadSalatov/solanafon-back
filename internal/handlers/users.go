package handlers

import (
	"fmt"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UsersHandler handles /api/users/* endpoints
type UsersHandler struct {
	db *gorm.DB
}

func NewUsersHandler(db *gorm.DB) *UsersHandler {
	return &UsersHandler{db: db}
}

// GetMe — GET /api/users/me
func (h *UsersHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "User not found"}})
	}

	var appsUsed, transactions, nfts, appsCreated int64
	h.db.Model(&models.AppUser{}).Where("user_id = ?", userID).Count(&appsUsed)
	h.db.Model(&models.ManaTransaction{}).Where("user_id = ?", userID).Count(&transactions)
	h.db.Model(&models.MiniApp{}).Where("creator_id = ?", userID).Count(&appsCreated)

	return c.JSON(fiber.Map{
		"id":          fmt.Sprintf("user_%d", user.ID),
		"email":       user.Email,
		"displayName": user.GetDisplayName(),
		"avatarUrl":   user.GetAvatarURL(),
		"createdAt":   user.CreatedAt,
		"stats": fiber.Map{
			"appsUsed":    appsUsed,
			"transactions": transactions,
			"nftsOwned":   nfts,
			"appsCreated": appsCreated,
		},
		"settings": fiber.Map{
			"language":           user.Language,
			"hapticFeedback":     user.HapticFeedback,
			"pushNotifications":  user.PushNotifications,
			"emailNotifications": user.EmailNotifications,
			"marketingEmails":    user.MarketingEmails,
			"biometricEnabled":   user.BiometricEnabled,
			"twoFactorEnabled":   user.TwoFactorEnabled,
		},
	})
}

// UpdateMe — PATCH /api/users/me
func (h *UsersHandler) UpdateMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "User not found"}})
	}

	var input struct {
		DisplayName *string `json:"displayName"`
		AvatarURL   *string `json:"avatarUrl"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "Invalid request body"}})
	}

	if input.DisplayName != nil {
		user.DisplayName = *input.DisplayName
		user.Name = *input.DisplayName
	}
	if input.AvatarURL != nil {
		user.AvatarURL = *input.AvatarURL
		user.Avatar = *input.AvatarURL
	}
	h.db.Save(&user)

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id": fmt.Sprintf("user_%d", user.ID), "email": user.Email,
			"displayName": user.GetDisplayName(), "avatarUrl": user.GetAvatarURL(),
		},
	})
}

// UpdateSettings — PATCH /api/users/me/settings
func (h *UsersHandler) UpdateSettings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "User not found"}})
	}

	var input struct {
		Language           *string `json:"language"`
		HapticFeedback     *bool   `json:"hapticFeedback"`
		PushNotifications  *bool   `json:"pushNotifications"`
		EmailNotifications *bool   `json:"emailNotifications"`
		MarketingEmails    *bool   `json:"marketingEmails"`
		BiometricEnabled   *bool   `json:"biometricEnabled"`
		TwoFactorEnabled   *bool   `json:"twoFactorEnabled"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "Invalid body"}})
	}

	if input.Language != nil {
		user.Language = *input.Language
	}
	if input.HapticFeedback != nil {
		user.HapticFeedback = *input.HapticFeedback
	}
	if input.PushNotifications != nil {
		user.PushNotifications = *input.PushNotifications
	}
	if input.EmailNotifications != nil {
		user.EmailNotifications = *input.EmailNotifications
	}
	if input.MarketingEmails != nil {
		user.MarketingEmails = *input.MarketingEmails
	}
	if input.BiometricEnabled != nil {
		user.BiometricEnabled = *input.BiometricEnabled
	}
	if input.TwoFactorEnabled != nil {
		user.TwoFactorEnabled = *input.TwoFactorEnabled
	}
	h.db.Save(&user)

	return c.JSON(fiber.Map{
		"success": true,
		"settings": fiber.Map{
			"language": user.Language, "hapticFeedback": user.HapticFeedback,
			"pushNotifications": user.PushNotifications, "emailNotifications": user.EmailNotifications,
			"marketingEmails": user.MarketingEmails, "biometricEnabled": user.BiometricEnabled,
			"twoFactorEnabled": user.TwoFactorEnabled,
		},
	})
}

// DeleteMe — DELETE /api/users/me
func (h *UsersHandler) DeleteMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	h.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{})
	h.db.Where("user_id = ?", userID).Delete(&models.Session{})
	h.db.Delete(&models.User{}, userID)
	return c.JSON(fiber.Map{"success": true})
}

// GetSessions — GET /api/users/me/sessions
func (h *UsersHandler) GetSessions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var sessions []models.Session
	h.db.Where("user_id = ?", userID).Order("last_active DESC").Find(&sessions)

	result := make([]fiber.Map, len(sessions))
	for i, s := range sessions {
		result[i] = fiber.Map{
			"id": fmt.Sprintf("sess_%d", s.ID), "deviceType": s.DeviceType,
			"deviceName": s.DeviceName, "ipAddress": s.IPAddress,
			"location": s.Location, "lastActive": s.LastActive, "isCurrent": false,
		}
	}
	return c.JSON(fiber.Map{"sessions": result})
}

// RevokeSession — DELETE /api/users/me/sessions/:sessionId
func (h *UsersHandler) RevokeSession(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	sessionID := c.Params("sessionId")
	h.db.Where("user_id = ? AND id = ?", userID, sessionID).Delete(&models.Session{})
	return c.JSON(fiber.Map{"success": true})
}

// RevokeAllSessions — DELETE /api/users/me/sessions
func (h *UsersHandler) RevokeAllSessions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	h.db.Where("user_id = ?", userID).Delete(&models.Session{})
	h.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{})
	return c.JSON(fiber.Map{"success": true})
}
