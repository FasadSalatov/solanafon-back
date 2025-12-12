package handlers

import (
	"strings"
	"time"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// BotHandler - handles Bot API requests from external services
type BotHandler struct {
	db *gorm.DB
}

func NewBotHandler(db *gorm.DB) *BotHandler {
	return &BotHandler{db: db}
}

// SendMessageInput - input for sending message via Bot API
type SendMessageInput struct {
	ChatID      uint   `json:"chat_id"`
	Text        string `json:"text"`
	MessageType string `json:"message_type,omitempty"` // text, image, button
	Metadata    string `json:"metadata,omitempty"`     // JSON for buttons, images, etc.
}

// SendMessage - send message to user from bot (requires API token)
// POST /bot/sendMessage
func (h *BotHandler) SendMessage(c *fiber.Ctx) error {
	// Get app from API token
	app, err := h.getAppFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":          false,
			"error_code":  401,
			"description": "Unauthorized: invalid API token",
		})
	}

	var input SendMessageInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":          false,
			"error_code":  400,
			"description": "Bad Request: invalid request body",
		})
	}

	if input.ChatID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":          false,
			"error_code":  400,
			"description": "Bad Request: chat_id is required",
		})
	}

	if strings.TrimSpace(input.Text) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":          false,
			"error_code":  400,
			"description": "Bad Request: text is required",
		})
	}

	// Verify user exists
	var user models.User
	if err := h.db.First(&user, input.ChatID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":          false,
			"error_code":  400,
			"description": "Bad Request: chat not found",
		})
	}

	// Create message
	messageType := "text"
	if input.MessageType != "" {
		messageType = input.MessageType
	}

	msg := models.AppMessage{
		AppID:       app.ID,
		UserID:      user.ID,
		Content:     input.Text,
		IsFromBot:   true,
		IsRead:      false,
		MessageType: messageType,
		Metadata:    input.Metadata,
		CreatedAt:   time.Now(),
	}

	if err := h.db.Create(&msg).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":          false,
			"error_code":  500,
			"description": "Internal Server Error: failed to send message",
		})
	}

	return c.JSON(fiber.Map{
		"ok": true,
		"result": fiber.Map{
			"message_id": msg.ID,
			"chat": fiber.Map{
				"id":   user.ID,
				"type": "private",
			},
			"date": msg.CreatedAt.Unix(),
			"text": msg.Content,
		},
	})
}

// GetUpdates - get pending messages for the bot (polling mode)
// GET /bot/getUpdates
func (h *BotHandler) GetUpdates(c *fiber.Ctx) error {
	app, err := h.getAppFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":          false,
			"error_code":  401,
			"description": "Unauthorized: invalid API token",
		})
	}

	// Get unprocessed user messages (not from bot)
	var messages []models.AppMessage
	if err := h.db.Preload("User").
		Where("app_id = ? AND is_from_bot = ? AND is_read = ?", app.ID, false, false).
		Order("created_at ASC").
		Limit(100).
		Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":          false,
			"error_code":  500,
			"description": "Internal Server Error: failed to fetch updates",
		})
	}

	// Mark as read
	if len(messages) > 0 {
		var ids []uint
		for _, m := range messages {
			ids = append(ids, m.ID)
		}
		h.db.Model(&models.AppMessage{}).Where("id IN ?", ids).Update("is_read", true)
	}

	// Format updates
	var updates []fiber.Map
	for _, msg := range messages {
		var user models.User
		h.db.First(&user, msg.UserID)

		updates = append(updates, fiber.Map{
			"update_id": msg.ID,
			"message": fiber.Map{
				"message_id": msg.ID,
				"from": fiber.Map{
					"id":       user.ID,
					"email":    user.Email,
					"name":     user.Name,
					"language": user.Language,
				},
				"chat": fiber.Map{
					"id":   user.ID,
					"type": "private",
				},
				"date": msg.CreatedAt.Unix(),
				"text": msg.Content,
			},
		})
	}

	return c.JSON(fiber.Map{
		"ok":     true,
		"result": updates,
	})
}

// SetWebhook - set webhook URL for the bot
// POST /bot/setWebhook
func (h *BotHandler) SetWebhook(c *fiber.Ctx) error {
	app, err := h.getAppFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":          false,
			"error_code":  401,
			"description": "Unauthorized: invalid API token",
		})
	}

	var input struct {
		URL string `json:"url"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":          false,
			"error_code":  400,
			"description": "Bad Request: invalid request body",
		})
	}

	app.WebhookURL = input.URL
	if err := h.db.Save(&app).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":          false,
			"error_code":  500,
			"description": "Internal Server Error: failed to set webhook",
		})
	}

	return c.JSON(fiber.Map{
		"ok":          true,
		"result":      true,
		"description": "Webhook was set",
	})
}

// DeleteWebhook - remove webhook URL
// POST /bot/deleteWebhook
func (h *BotHandler) DeleteWebhook(c *fiber.Ctx) error {
	app, err := h.getAppFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":          false,
			"error_code":  401,
			"description": "Unauthorized: invalid API token",
		})
	}

	app.WebhookURL = ""
	if err := h.db.Save(&app).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":          false,
			"error_code":  500,
			"description": "Internal Server Error: failed to delete webhook",
		})
	}

	return c.JSON(fiber.Map{
		"ok":          true,
		"result":      true,
		"description": "Webhook was deleted",
	})
}

// GetWebhookInfo - get current webhook info
// GET /bot/getWebhookInfo
func (h *BotHandler) GetWebhookInfo(c *fiber.Ctx) error {
	app, err := h.getAppFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":          false,
			"error_code":  401,
			"description": "Unauthorized: invalid API token",
		})
	}

	// Get pending updates count
	var pendingCount int64
	h.db.Model(&models.AppMessage{}).
		Where("app_id = ? AND is_from_bot = ? AND is_read = ?", app.ID, false, false).
		Count(&pendingCount)

	return c.JSON(fiber.Map{
		"ok": true,
		"result": fiber.Map{
			"url":                  app.WebhookURL,
			"has_custom_certificate": false,
			"pending_update_count": pendingCount,
		},
	})
}

// GetMe - get bot info
// GET /bot/getMe
func (h *BotHandler) GetMe(c *fiber.Ctx) error {
	app, err := h.getAppFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":          false,
			"error_code":  401,
			"description": "Unauthorized: invalid API token",
		})
	}

	return c.JSON(fiber.Map{
		"ok": true,
		"result": fiber.Map{
			"id":         app.ID,
			"is_bot":     true,
			"title":      app.Title,
			"username":   app.BotUsername,
			"can_join_groups": false,
			"can_read_all_group_messages": false,
			"supports_inline_queries": false,
		},
	})
}

// SetCommands - set bot commands
// POST /bot/setMyCommands
func (h *BotHandler) SetCommands(c *fiber.Ctx) error {
	app, err := h.getAppFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":          false,
			"error_code":  401,
			"description": "Unauthorized: invalid API token",
		})
	}

	var input struct {
		Commands []struct {
			Command     string `json:"command"`
			Description string `json:"description"`
		} `json:"commands"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":          false,
			"error_code":  400,
			"description": "Bad Request: invalid request body",
		})
	}

	// Delete existing commands
	h.db.Where("app_id = ?", app.ID).Delete(&models.BotCommand{})

	// Create new commands
	for _, cmd := range input.Commands {
		command := strings.TrimPrefix(cmd.Command, "/")
		botCmd := models.BotCommand{
			AppID:       app.ID,
			Command:     "/" + command,
			Description: cmd.Description,
			IsEnabled:   true,
		}
		h.db.Create(&botCmd)
	}

	return c.JSON(fiber.Map{
		"ok":     true,
		"result": true,
	})
}

// GetCommands - get bot commands
// GET /bot/getMyCommands
func (h *BotHandler) GetCommands(c *fiber.Ctx) error {
	app, err := h.getAppFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":          false,
			"error_code":  401,
			"description": "Unauthorized: invalid API token",
		})
	}

	var commands []models.BotCommand
	h.db.Where("app_id = ? AND is_enabled = ?", app.ID, true).Find(&commands)

	var result []fiber.Map
	for _, cmd := range commands {
		result = append(result, fiber.Map{
			"command":     strings.TrimPrefix(cmd.Command, "/"),
			"description": cmd.Description,
		})
	}

	return c.JSON(fiber.Map{
		"ok":     true,
		"result": result,
	})
}

// getAppFromToken - extract app from API token in Authorization header
func (h *BotHandler) getAppFromToken(c *fiber.Ctx) (*models.MiniApp, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, fiber.ErrUnauthorized
	}

	// Support both "Bearer <token>" and just "<token>"
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		return nil, fiber.ErrUnauthorized
	}

	var app models.MiniApp
	if err := h.db.Where("api_token = ?", token).First(&app).Error; err != nil {
		return nil, fiber.ErrUnauthorized
	}

	// Check if app is approved
	if app.ModerationStatus != models.ModerationApproved {
		return nil, fiber.ErrUnauthorized
	}

	return &app, nil
}
