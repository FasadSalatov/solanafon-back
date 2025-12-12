package handlers

import (
	"time"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SecretHandler struct {
	db *gorm.DB
}

func NewSecretHandler(db *gorm.DB) *SecretHandler {
	return &SecretHandler{db: db}
}

// GetNumbers - get available secret numbers
func (h *SecretHandler) GetNumbers(c *fiber.Ctx) error {
	var numbers []models.SecretNumber
	if err := h.db.Where("is_available = ?", true).
		Order("is_premium ASC, price_mp ASC").
		Find(&numbers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch numbers",
		})
	}

	// Get user's MP balance for frontend
	userID := c.Locals("userID").(uint)
	var user models.User
	h.db.First(&user, userID)

	return c.JSON(fiber.Map{
		"numbers":        numbers,
		"total":          len(numbers),
		"userManaPoints": user.ManaPoints,
	})
}

// ActivateNumberInput - input for activating secret number
type ActivateNumberInput struct {
	NumberID uint `json:"numberId"`
}

// ActivateNumber - purchase and activate a secret number with MP
func (h *SecretHandler) ActivateNumber(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input ActivateNumberInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if user already has active secret access
	var existingAccess models.SecretAccess
	if err := h.db.Where("user_id = ? AND is_active = ?", userID, true).First(&existingAccess).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "You already have active secret access",
		})
	}

	// Check if number exists and is available
	var number models.SecretNumber
	if err := h.db.First(&number, input.NumberID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Number not found",
		})
	}

	if !number.IsAvailable {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Number is not available",
		})
	}

	// Check if user has enough MP
	if user.ManaPoints < number.PriceMP {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":    "Insufficient Mana Points",
			"required": number.PriceMP,
			"balance":  user.ManaPoints,
		})
	}

	// Start transaction
	tx := h.db.Begin()

	// Deduct MP from user
	user.ManaPoints -= number.PriceMP
	user.HasSecretAccess = true
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Create MP transaction record
	mpTransaction := models.ManaTransaction{
		UserID:      userID,
		Amount:      -number.PriceMP, // negative = debit
		Type:        "purchase",
		Description: "Secret number activation: " + number.Number,
		ReferenceID: &number.ID,
		CreatedAt:   time.Now(),
	}
	if err := tx.Create(&mpTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create transaction",
		})
	}

	// Create secret access
	secretAccess := models.SecretAccess{
		UserID:         userID,
		SecretNumberID: number.ID,
		ActivatedAt:    time.Now(),
		IsActive:       true,
	}
	if err := tx.Create(&secretAccess).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to activate secret access",
		})
	}

	// Mark number as unavailable
	number.IsAvailable = false
	if err := tx.Save(&number).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update number",
		})
	}

	// Commit transaction
	tx.Commit()

	// Load full data for response
	h.db.Preload("SecretNumber").First(&secretAccess, secretAccess.ID)

	return c.JSON(fiber.Map{
		"message":        "Secret access activated successfully",
		"secretAccess":   secretAccess,
		"newManaBalance": user.ManaPoints,
		"spent":          number.PriceMP,
	})
}

// GetStatus - get current secret access status
func (h *SecretHandler) GetStatus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var secretAccess models.SecretAccess
	hasAccess := h.db.Where("user_id = ? AND is_active = ?", userID, true).
		Preload("SecretNumber").
		First(&secretAccess).Error == nil

	response := fiber.Map{
		"hasAccess":  user.HasSecretAccess,
		"manaPoints": user.ManaPoints,
	}

	if hasAccess {
		response["activatedAt"] = secretAccess.ActivatedAt
		response["number"] = secretAccess.SecretNumber.Number
		response["isPremium"] = secretAccess.SecretNumber.IsPremium
	}

	// Show benefits of secret access
	response["benefits"] = []string{
		"Полная анонимность и конфиденциальность",
		"Получение одноразовых SMS кодов",
		"Без привязки к реальному номеру телефона",
		"Доступ к секретным приложениям",
	}

	return c.JSON(response)
}

// DeactivateSecret - deactivate secret access (admin or user request)
func (h *SecretHandler) DeactivateSecret(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var secretAccess models.SecretAccess
	if err := h.db.Where("user_id = ? AND is_active = ?", userID, true).First(&secretAccess).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No active secret access found",
		})
	}

	// Deactivate
	secretAccess.IsActive = false
	h.db.Save(&secretAccess)

	// Update user
	var user models.User
	h.db.First(&user, userID)
	user.HasSecretAccess = false
	h.db.Save(&user)

	// Make number available again
	h.db.Model(&models.SecretNumber{}).
		Where("id = ?", secretAccess.SecretNumberID).
		Update("is_available", true)

	return c.JSON(fiber.Map{
		"message": "Secret access deactivated",
	})
}
