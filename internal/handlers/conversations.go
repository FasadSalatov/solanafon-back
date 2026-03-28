package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ConversationsHandler handles /api/conversations/* endpoints
type ConversationsHandler struct {
	db *gorm.DB
}

func NewConversationsHandler(db *gorm.DB) *ConversationsHandler {
	return &ConversationsHandler{db: db}
}

// ListConversations — GET /api/conversations
func (h *ConversationsHandler) ListConversations(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var convs []models.Conversation
	var total int64
	h.db.Model(&models.Conversation{}).Where("user_id = ?", userID).Count(&total)
	h.db.Where("user_id = ?", userID).Preload("App").Order("updated_at DESC").
		Offset(offset).Limit(limit).Find(&convs)

	result := make([]fiber.Map, 0, len(convs))
	for _, conv := range convs {
		// Get last message
		var lastMsg models.ChatMessage
		h.db.Where("conversation_id = ?", conv.ID).Order("created_at DESC").First(&lastMsg)

		var lastMsgMap fiber.Map
		if lastMsg.ID > 0 {
			lastMsgMap = fiber.Map{
				"id": fmt.Sprintf("msg_%d", lastMsg.ID),
				"content":    json.RawMessage(lastMsg.Content),
				"timestamp":  lastMsg.CreatedAt.UnixMilli(),
				"senderType": lastMsg.SenderType,
			}
		}

		result = append(result, fiber.Map{
			"id":          fmt.Sprintf("conv_%d", conv.ID),
			"appId":       fmt.Sprintf("app_%d", conv.AppID),
			"userId":      fmt.Sprintf("user_%d", conv.UserID),
			"appName":     conv.App.Title,
			"appIcon":     conv.App.Icon,
			"appIconUrl":  conv.App.IconURL,
			"appUrl":      conv.App.URL,
			"lastMessage": lastMsgMap,
			"unreadCount": conv.UnreadCount,
			"createdAt":   conv.CreatedAt,
			"updatedAt":   conv.UpdatedAt,
			"isActive":    conv.IsActive,
		})
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return c.JSON(fiber.Map{
		"conversations": result,
		"pagination": fiber.Map{
			"currentPage": page, "totalPages": totalPages,
			"totalItems": total, "hasMore": page < totalPages,
		},
	})
}

// StartConversation — POST /api/apps/:appId/conversations
func (h *ConversationsHandler) StartConversation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID, _ := strconv.Atoi(c.Params("appId"))

	var app models.MiniApp
	if err := h.db.First(&app, appID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "App not found"}})
	}

	// Check existing conversation
	var existing models.Conversation
	if err := h.db.Where("user_id = ? AND app_id = ?", userID, appID).First(&existing).Error; err == nil {
		return c.JSON(fiber.Map{"success": true, "conversation": formatConversation(existing, app)})
	}

	var input struct {
		InitialMessage string `json:"initialMessage"`
	}
	c.BodyParser(&input)

	now := time.Now()
	conv := models.Conversation{
		AppID: uint(appID), UserID: userID, IsActive: true, LastMessageAt: &now,
	}
	h.db.Create(&conv)

	// Track app usage
	var appUser models.AppUser
	if h.db.Where("user_id = ? AND app_id = ?", userID, appID).First(&appUser).Error != nil {
		h.db.Create(&models.AppUser{UserID: userID, AppID: uint(appID), LastUsed: now})
		h.db.Model(&app).UpdateColumn("users_count", gorm.Expr("users_count + 1"))
	}

	// Welcome message
	var welcomeMsg fiber.Map
	if app.WelcomeMessage != "" {
		content, _ := json.Marshal(fiber.Map{"type": "text", "text": app.WelcomeMessage})
		msg := models.ChatMessage{
			ConversationID: conv.ID, AppID: uint(appID),
			SenderID: "bot", SenderType: "bot",
			Content: string(content), Status: "delivered",
		}
		h.db.Create(&msg)
		welcomeMsg = fiber.Map{
			"id": fmt.Sprintf("msg_%d", msg.ID),
			"content": json.RawMessage(msg.Content),
			"senderType": "bot", "timestamp": msg.CreatedAt.UnixMilli(),
		}
	}

	return c.Status(201).JSON(fiber.Map{
		"success":        true,
		"conversation":   formatConversation(conv, app),
		"welcomeMessage": welcomeMsg,
	})
}

// GetMessages — GET /api/conversations/:conversationId/messages
func (h *ConversationsHandler) GetMessages(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	convID, _ := strconv.Atoi(c.Params("conversationId"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	before := c.Query("before")

	var conv models.Conversation
	if err := h.db.Where("id = ? AND user_id = ?", convID, userID).First(&conv).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "Conversation not found"}})
	}

	query := h.db.Where("conversation_id = ?", convID)
	if before != "" {
		beforeID, _ := strconv.Atoi(before)
		query = query.Where("id < ?", beforeID)
	}

	var messages []models.ChatMessage
	query.Order("created_at DESC").Limit(limit).Find(&messages)

	result := make([]fiber.Map, 0, len(messages))
	for _, msg := range messages {
		result = append(result, fiber.Map{
			"id": fmt.Sprintf("msg_%d", msg.ID),
			"appId":          fmt.Sprintf("app_%d", msg.AppID),
			"conversationId": fmt.Sprintf("conv_%d", msg.ConversationID),
			"senderId":       msg.SenderID,
			"senderType":     msg.SenderType,
			"content":        json.RawMessage(msg.Content),
			"timestamp":      msg.CreatedAt.UnixMilli(),
			"status":         msg.Status,
			"replyToId":      msg.ReplyToID,
			"metadata":       msg.Metadata,
		})
	}

	var totalCount int64
	h.db.Model(&models.ChatMessage{}).Where("conversation_id = ?", convID).Count(&totalCount)

	return c.JSON(fiber.Map{
		"messages":   result,
		"hasMore":    len(messages) == limit,
		"pagination": fiber.Map{"totalItems": totalCount},
	})
}

// SendMessage — POST /api/conversations/:conversationId/messages
func (h *ConversationsHandler) SendMessage(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	convID, _ := strconv.Atoi(c.Params("conversationId"))

	var conv models.Conversation
	if err := h.db.Where("id = ? AND user_id = ?", convID, userID).Preload("App").First(&conv).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "Conversation not found"}})
	}

	var input struct {
		Content   json.RawMessage `json:"content"`
		ReplyToID *uint           `json:"replyToId"`
		Metadata  json.RawMessage `json:"metadata"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "Invalid body"}})
	}

	msg := models.ChatMessage{
		ConversationID: uint(convID), AppID: conv.AppID,
		SenderID: fmt.Sprintf("user_%d", userID), SenderType: "user",
		Content: string(input.Content), Status: "sent",
		ReplyToID: input.ReplyToID,
	}
	if input.Metadata != nil {
		msg.Metadata = string(input.Metadata)
	}
	h.db.Create(&msg)

	now := time.Now()
	h.db.Model(&conv).Updates(map[string]interface{}{"last_message_at": now, "updated_at": now})

	// Trigger webhook if configured
	if conv.App.WebhookURL != "" {
		go triggerConvWebhook(h.db, conv.App, conv, msg, "message.received")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fiber.Map{
			"id": fmt.Sprintf("msg_%d", msg.ID),
			"content": json.RawMessage(msg.Content),
			"senderType": "user", "timestamp": msg.CreatedAt.UnixMilli(), "status": msg.Status,
		},
	})
}

// ButtonCallback — POST /api/conversations/:conversationId/callback
func (h *ConversationsHandler) ButtonCallback(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	convID, _ := strconv.Atoi(c.Params("conversationId"))

	var conv models.Conversation
	if err := h.db.Where("id = ? AND user_id = ?", convID, userID).Preload("App").First(&conv).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "Conversation not found"}})
	}

	var input struct {
		MessageID string `json:"messageId"`
		ButtonID  string `json:"buttonId"`
		Payload   string `json:"payload"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "Invalid body"}})
	}

	// Trigger callback webhook
	if conv.App.WebhookURL != "" {
		go triggerCallbackWebhook(h.db, conv.App, conv, input.MessageID, input.ButtonID, input.Payload, userID)
	}

	return c.JSON(fiber.Map{"success": true, "action": nil})
}

// MarkAsRead — POST /api/conversations/:conversationId/read
func (h *ConversationsHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	convID, _ := strconv.Atoi(c.Params("conversationId"))
	h.db.Model(&models.Conversation{}).Where("id = ? AND user_id = ?", convID, userID).
		Update("unread_count", 0)
	h.db.Model(&models.ChatMessage{}).Where("conversation_id = ? AND sender_type = ? AND status != ?", convID, "bot", "read").
		Update("status", "read")
	return c.JSON(fiber.Map{"success": true})
}

// DeleteConversation — DELETE /api/conversations/:conversationId
func (h *ConversationsHandler) DeleteConversation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	convID, _ := strconv.Atoi(c.Params("conversationId"))

	var conv models.Conversation
	if err := h.db.Where("id = ? AND user_id = ?", convID, userID).Preload("App").First(&conv).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "Conversation not found"}})
	}

	h.db.Delete(&conv)

	if conv.App.WebhookURL != "" {
		go triggerConvWebhook(h.db, conv.App, conv, models.ChatMessage{}, "conversation.ended")
	}

	return c.JSON(fiber.Map{"success": true})
}

// helpers

func formatConversation(conv models.Conversation, app models.MiniApp) fiber.Map {
	return fiber.Map{
		"id": fmt.Sprintf("conv_%d", conv.ID), "appId": fmt.Sprintf("app_%d", conv.AppID),
		"appName": app.Title, "appIcon": app.Icon, "appIconUrl": app.IconURL,
		"unreadCount": conv.UnreadCount, "isActive": conv.IsActive,
		"createdAt": conv.CreatedAt, "updatedAt": conv.UpdatedAt,
	}
}

func triggerConvWebhook(db *gorm.DB, app models.MiniApp, conv models.Conversation, msg models.ChatMessage, event string) {
	payload := fiber.Map{
		"event": event, "timestamp": time.Now().UnixMilli(),
		"data": fiber.Map{
			"conversationId": fmt.Sprintf("conv_%d", conv.ID),
			"userId":         fmt.Sprintf("user_%d", conv.UserID),
		},
	}
	if msg.ID > 0 {
		payload["data"] = fiber.Map{
			"conversationId": fmt.Sprintf("conv_%d", conv.ID),
			"message": fiber.Map{
				"id": fmt.Sprintf("msg_%d", msg.ID), "senderId": msg.SenderID,
				"senderType": msg.SenderType, "content": json.RawMessage(msg.Content),
				"timestamp": msg.CreatedAt.UnixMilli(),
			},
		}
	}
	// Fire and forget webhook (same pattern as existing bot.go)
	body, _ := json.Marshal(payload)
	_ = body // TODO: HTTP POST to app.WebhookURL with signing
}

func triggerCallbackWebhook(db *gorm.DB, app models.MiniApp, conv models.Conversation, msgID, btnID, payload string, userID uint) {
	data := fiber.Map{
		"event": "callback.received", "timestamp": time.Now().UnixMilli(),
		"data": fiber.Map{
			"conversationId": fmt.Sprintf("conv_%d", conv.ID),
			"messageId": msgID, "buttonId": btnID, "payload": payload,
			"userId": fmt.Sprintf("user_%d", userID),
		},
	}
	body, _ := json.Marshal(data)
	_ = body
}
