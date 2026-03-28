package handlers

import (
	"fmt"
	"time"

	"github.com/fasad/solanafon-back/internal/config"
	"github.com/fasad/solanafon-back/internal/models"
	"github.com/fasad/solanafon-back/internal/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AuthV2Handler handles new /api/auth/* endpoints for the mobile app
type AuthV2Handler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthV2Handler(db *gorm.DB, cfg *config.Config) *AuthV2Handler {
	return &AuthV2Handler{db: db, cfg: cfg}
}

// SendCode — POST /api/auth/send-code
func (h *AuthV2Handler) SendCode(c *fiber.Ctx) error {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&input); err != nil || input.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "Email is required"}})
	}

	code, err := utils.GenerateOTP(6)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fiber.Map{"code": "INTERNAL_ERROR", "message": "Failed to generate code"}})
	}
	expiresAt := time.Now().Add(time.Duration(h.cfg.OTPExpiryMinutes) * time.Minute)

	otp := models.OTP{Email: input.Email, Code: code, ExpiresAt: expiresAt}
	h.db.Create(&otp)

	emailCfg := utils.EmailConfig{
		Host: h.cfg.SMTPHost, Port: h.cfg.SMTPPort,
		User: h.cfg.SMTPUser, Password: h.cfg.SMTPPassword,
	}
	if err := utils.SendOTPEmail(input.Email, code, emailCfg); err != nil {
		fmt.Printf("[DEV] OTP for %s: %s\n", input.Email, code)
	}

	return c.JSON(fiber.Map{"success": true, "message": "Code sent to email", "expiresIn": h.cfg.OTPExpiryMinutes * 60})
}

// VerifyCode — POST /api/auth/verify-code
func (h *AuthV2Handler) VerifyCode(c *fiber.Ctx) error {
	var input struct {
		Email        string `json:"email"`
		Code         string `json:"code"`
		ReferralCode string `json:"referralCode"`
	}
	if err := c.BodyParser(&input); err != nil || input.Email == "" || input.Code == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "Email and code are required"}})
	}

	var otp models.OTP
	if err := h.db.Where("email = ? AND verified = false AND expires_at > ?", input.Email, time.Now()).
		Order("created_at DESC").First(&otp).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "INVALID_CODE", "message": "Invalid or expired code"}})
	}

	if otp.Code != input.Code {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "INVALID_CODE", "message": "Invalid code"}})
	}

	otp.Verified = true
	h.db.Save(&otp)

	var user models.User
	isNewUser := false
	if err := h.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		isNewUser = true
		user = models.User{
			Email:        input.Email,
			Name:         "Solana User",
			DisplayName:  "Solana User",
			ManaPoints:   500,
			ReferralCode: models.GenerateReferralCode(),
		}
		h.db.Create(&user)

		h.db.Create(&models.ManaTransaction{
			UserID: user.ID, Amount: 500, Type: "reward", Description: "Welcome bonus",
		})

		// Process referral
		if input.ReferralCode != "" {
			var referrer models.User
			if err := h.db.Where("referral_code = ?", input.ReferralCode).First(&referrer).Error; err == nil {
				h.db.Create(&models.Referral{
					ReferrerID: referrer.ID, ReferredID: user.ID, RegisteredAt: time.Now(),
				})
				// Bonus for referrer
				referrer.ManaPoints += 100
				h.db.Save(&referrer)
				h.db.Create(&models.ManaTransaction{
					UserID: referrer.ID, Amount: 100, Type: "reward", Description: "Referral bonus",
				})
			}
		}
	}

	// Generate tokens
	token, _ := utils.GenerateJWT(user.ID, user.Email, h.cfg.JWTSecret)
	refreshTokenStr := models.GenerateRefreshToken()

	h.db.Create(&models.RefreshToken{
		UserID: user.ID, Token: refreshTokenStr, ExpiresAt: time.Now().Add(90 * 24 * time.Hour),
	})

	// Create session
	h.db.Create(&models.Session{
		UserID: user.ID, DeviceType: "android", IPAddress: c.IP(), LastActive: time.Now(),
	})

	_ = isNewUser
	return c.JSON(fiber.Map{
		"success":      true,
		"token":        token,
		"refreshToken": refreshTokenStr,
		"user": fiber.Map{
			"id":          fmt.Sprintf("user_%d", user.ID),
			"email":       user.Email,
			"displayName": user.GetDisplayName(),
			"avatarUrl":   user.GetAvatarURL(),
		},
	})
}

// RefreshToken — POST /api/auth/refresh
func (h *AuthV2Handler) RefreshToken(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.BodyParser(&input); err != nil || input.RefreshToken == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "refreshToken is required"}})
	}

	var rt models.RefreshToken
	if err := h.db.Where("token = ? AND expires_at > ?", input.RefreshToken, time.Now()).First(&rt).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": fiber.Map{"code": "INVALID_TOKEN", "message": "Invalid or expired refresh token"}})
	}

	var user models.User
	h.db.First(&user, rt.UserID)

	// Delete old, create new
	h.db.Delete(&rt)

	newToken, _ := utils.GenerateJWT(user.ID, user.Email, h.cfg.JWTSecret)
	newRefresh := models.GenerateRefreshToken()
	h.db.Create(&models.RefreshToken{
		UserID: user.ID, Token: newRefresh, ExpiresAt: time.Now().Add(90 * 24 * time.Hour),
	})

	return c.JSON(fiber.Map{"token": newToken, "refreshToken": newRefresh})
}

// Logout — POST /api/logout
func (h *AuthV2Handler) LogoutV2(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	// Delete all refresh tokens for user
	h.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{})

	return c.JSON(fiber.Map{"success": true})
}
