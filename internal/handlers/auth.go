package handlers

import (
	"log"
	"net/mail"
	"time"

	"github.com/fasad/solanafon-back/internal/config"
	"github.com/fasad/solanafon-back/internal/models"
	"github.com/fasad/solanafon-back/internal/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	// Initial MP bonus for new users
	NewUserManaBonus = 500
)

type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

type EmailRequestInput struct {
	Email string `json:"email"`
}

type EmailVerifyInput struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// RequestOTP - send OTP code to email
func (h *AuthHandler) RequestOTP(c *fiber.Ctx) error {
	var input EmailRequestInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate email
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email address",
		})
	}

	// Generate OTP
	code, err := utils.GenerateOTP(6)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate OTP",
		})
	}

	// Delete old unverified OTPs for this email
	h.db.Where("email = ? AND verified = ?", input.Email, false).Delete(&models.OTP{})

	// Save OTP to database
	otp := models.OTP{
		Email:     input.Email,
		Code:      code,
		ExpiresAt: time.Now().Add(time.Duration(h.cfg.OTPExpiryMinutes) * time.Minute),
		Verified:  false,
	}

	if err := h.db.Create(&otp).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save OTP",
		})
	}

	// Try to send email (don't fail if SMTP not configured)
	emailConfig := utils.EmailConfig{
		Host:     h.cfg.SMTPHost,
		Port:     h.cfg.SMTPPort,
		User:     h.cfg.SMTPUser,
		Password: h.cfg.SMTPPassword,
	}

	if h.cfg.SMTPUser != "" && h.cfg.SMTPPassword != "" {
		if err := utils.SendOTPEmail(input.Email, code, emailConfig); err != nil {
			log.Printf("Failed to send email: %v", err)
			// Continue without failing - user can see code in logs during development
		}
	} else {
		// Development mode - log the code
		log.Printf("📧 OTP Code for %s: %s", input.Email, code)
	}

	return c.JSON(fiber.Map{
		"message": "OTP sent to your email",
		"email":   input.Email,
	})
}

// VerifyOTP - verify OTP code and return JWT token
func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var input EmailVerifyInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Find OTP
	var otp models.OTP
	if err := h.db.Where("email = ? AND code = ? AND verified = ?", input.Email, input.Code, false).
		Order("created_at DESC").
		First(&otp).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or expired OTP",
		})
	}

	// Check if expired
	if time.Now().After(otp.ExpiresAt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "OTP has expired",
		})
	}

	// Mark OTP as verified
	otp.Verified = true
	h.db.Save(&otp)

	// Find or create user
	var user models.User
	isNewUser := false
	result := h.db.Where("email = ?", input.Email).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		isNewUser = true
		// Create new user with initial Mana Points bonus
		user = models.User{
			Email:                input.Email,
			Name:                 "Solana User",
			ManaPoints:           NewUserManaBonus,
			HasSecretAccess:      false,
			Language:             "en",
			NotificationsEnabled: true,
			TwoFactorEnabled:     false,
		}
		if err := h.db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}

		// Create welcome MP transaction
		welcomeTransaction := models.ManaTransaction{
			UserID:      user.ID,
			Amount:      NewUserManaBonus,
			Type:        "reward",
			Description: "Welcome bonus for new user",
			CreatedAt:   time.Now(),
		}
		h.db.Create(&welcomeTransaction)
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, h.cfg.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	response := fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":              user.ID,
			"email":           user.Email,
			"name":            user.Name,
			"manaPoints":      user.ManaPoints,
			"hasSecretAccess": user.HasSecretAccess,
			"language":        user.Language,
			"createdAt":       user.CreatedAt,
		},
		"isNewUser": isNewUser,
	}

	if isNewUser {
		response["welcomeBonus"] = NewUserManaBonus
	}

	return c.JSON(response)
}

// GetMe - get current authenticated user
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get user stats
	var appsCount int64
	h.db.Model(&models.MiniApp{}).Where("creator_id = ?", userID).Count(&appsCount)

	var transactionsCount int64
	h.db.Model(&models.ManaTransaction{}).Where("user_id = ?", userID).Count(&transactionsCount)

	return c.JSON(fiber.Map{
		"id":                   user.ID,
		"email":                user.Email,
		"name":                 user.Name,
		"avatar":               user.Avatar,
		"manaPoints":           user.ManaPoints,
		"hasSecretAccess":      user.HasSecretAccess,
		"language":             user.Language,
		"notificationsEnabled": user.NotificationsEnabled,
		"twoFactorEnabled":     user.TwoFactorEnabled,
		"stats": fiber.Map{
			"appsCount":         appsCount,
			"transactionsCount": transactionsCount,
			"nftsCount":         0, // Placeholder
		},
		"createdAt": user.CreatedAt,
	})
}

// Logout - invalidate token (client-side implementation)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// JWT tokens are stateless, so logout is handled client-side
	// This endpoint just confirms the logout intention
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
