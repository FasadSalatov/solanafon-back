package handlers

import (
	"time"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProfileHandler struct {
	db *gorm.DB
}

func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

// GetProfile - get full user profile with stats
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get stats
	stats := h.getUserStats(userID)

	// Get Mana Points info
	manaInfo := fiber.Map{
		"balance":  user.ManaPoints,
		"valueUSD": float64(user.ManaPoints) * 0.023, // Example rate: 1 MP = $0.023
	}

	// Check secret access
	var secretAccess models.SecretAccess
	hasActiveSecret := h.db.Where("user_id = ? AND is_active = ?", userID, true).
		Preload("SecretNumber").
		First(&secretAccess).Error == nil

	response := fiber.Map{
		"id":                   user.ID,
		"email":                user.Email,
		"name":                 user.Name,
		"avatar":               user.Avatar,
		"manaPoints":           manaInfo,
		"hasSecretAccess":      user.HasSecretAccess,
		"language":             user.Language,
		"notificationsEnabled": user.NotificationsEnabled,
		"twoFactorEnabled":     user.TwoFactorEnabled,
		"stats":                stats,
		"createdAt":            user.CreatedAt,
	}

	if hasActiveSecret {
		response["secretAccess"] = fiber.Map{
			"activatedAt": secretAccess.ActivatedAt,
			"number":      secretAccess.SecretNumber.Number,
			"isPremium":   secretAccess.SecretNumber.IsPremium,
		}
	}

	return c.JSON(response)
}

// getUserStats - get computed stats for user
func (h *ProfileHandler) getUserStats(userID uint) models.UserStats {
	var stats models.UserStats

	// Count user's created apps
	h.db.Model(&models.MiniApp{}).Where("creator_id = ?", userID).Count((*int64)(&stats.AppsCount))

	// Count mana transactions
	h.db.Model(&models.ManaTransaction{}).Where("user_id = ?", userID).Count((*int64)(&stats.TransactionsCount))

	// NFTs count (placeholder - would come from blockchain)
	stats.NFTsCount = 0

	return stats
}

// UpdateProfileInput - input for profile update
type UpdateProfileInput struct {
	Name     string `json:"name,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Language string `json:"language,omitempty"`
}

// UpdateProfile - update user profile
func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input UpdateProfileInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Update fields if provided
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Avatar != "" {
		user.Avatar = input.Avatar
	}
	if input.Language != "" {
		user.Language = input.Language
	}

	if err := h.db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user": fiber.Map{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"avatar":   user.Avatar,
			"language": user.Language,
		},
	})
}

// UpdateSettingsInput - input for settings update
type UpdateSettingsInput struct {
	NotificationsEnabled *bool `json:"notificationsEnabled,omitempty"`
	TwoFactorEnabled     *bool `json:"twoFactorEnabled,omitempty"`
}

// UpdateSettings - update user settings
func (h *ProfileHandler) UpdateSettings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input UpdateSettingsInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if input.NotificationsEnabled != nil {
		user.NotificationsEnabled = *input.NotificationsEnabled
	}
	if input.TwoFactorEnabled != nil {
		user.TwoFactorEnabled = *input.TwoFactorEnabled
	}

	if err := h.db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update settings",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Settings updated successfully",
		"settings": fiber.Map{
			"notificationsEnabled": user.NotificationsEnabled,
			"twoFactorEnabled":     user.TwoFactorEnabled,
		},
	})
}

// GetManaHistory - get Mana Points transaction history
func (h *ProfileHandler) GetManaHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var transactions []models.ManaTransaction
	if err := h.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(50).
		Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}

	var user models.User
	h.db.First(&user, userID)

	return c.JSON(fiber.Map{
		"balance":      user.ManaPoints,
		"transactions": transactions,
		"total":        len(transactions),
	})
}

// TopUpManaInput - input for MP top up
type TopUpManaInput struct {
	Amount int `json:"amount"`
}

// TopUpMana - top up Mana Points (mock - would integrate with payment)
func (h *ProfileHandler) TopUpMana(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input TopUpManaInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be positive",
		})
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Create transaction
	transaction := models.ManaTransaction{
		UserID:      userID,
		Amount:      input.Amount,
		Type:        "topup",
		Description: "Mana Points top up",
		CreatedAt:   time.Now(),
	}

	if err := h.db.Create(&transaction).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create transaction",
		})
	}

	// Update balance
	user.ManaPoints += input.Amount
	h.db.Save(&user)

	return c.JSON(fiber.Map{
		"message":    "Mana Points added successfully",
		"newBalance": user.ManaPoints,
		"added":      input.Amount,
	})
}
