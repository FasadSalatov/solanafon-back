package handlers

import (
	"fmt"
	"strconv"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// NotificationsHandler handles /api/notifications/* endpoints
type NotificationsHandler struct {
	db *gorm.DB
}

func NewNotificationsHandler(db *gorm.DB) *NotificationsHandler {
	return &NotificationsHandler{db: db}
}

// RegisterPushToken — POST /api/notifications/register
func (h *NotificationsHandler) RegisterPushToken(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var input struct {
		FCMToken string `json:"fcmToken"`
		DeviceID string `json:"deviceId"`
		Platform string `json:"platform"`
	}
	if err := c.BodyParser(&input); err != nil || input.FCMToken == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "fcmToken is required"}})
	}

	// Upsert by device ID
	var existing models.PushToken
	if h.db.Where("user_id = ? AND device_id = ?", userID, input.DeviceID).First(&existing).Error == nil {
		existing.FCMToken = input.FCMToken
		h.db.Save(&existing)
	} else {
		platform := input.Platform
		if platform == "" {
			platform = "android"
		}
		h.db.Create(&models.PushToken{
			UserID: userID, FCMToken: input.FCMToken, DeviceID: input.DeviceID, Platform: platform,
		})
	}

	return c.JSON(fiber.Map{"success": true})
}

// UnregisterPushToken — DELETE /api/notifications/unregister
func (h *NotificationsHandler) UnregisterPushToken(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	h.db.Where("user_id = ?", userID).Delete(&models.PushToken{})
	return c.JSON(fiber.Map{"success": true})
}

// ListNotifications — GET /api/notifications
func (h *NotificationsHandler) ListNotifications(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	unreadOnly := c.Query("unreadOnly") == "true"
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	query := h.db.Where("user_id = ?", userID)
	if unreadOnly {
		query = query.Where("is_read = false")
	}

	var total int64
	query.Model(&models.Notification{}).Count(&total)

	var notifs []models.Notification
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&notifs)

	result := make([]fiber.Map, len(notifs))
	for i, n := range notifs {
		result[i] = fiber.Map{
			"id": fmt.Sprintf("notif_%d", n.ID), "title": n.Title,
			"body": n.Body, "type": n.Type, "isRead": n.IsRead,
			"createdAt": n.CreatedAt, "actionUrl": n.ActionURL,
		}
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"notifications": result,
		"pagination": fiber.Map{
			"page": page, "limit": limit, "total": total, "hasMore": page < totalPages,
		},
	})
}

// MarkAsRead — POST /api/notifications/:id/read
func (h *NotificationsHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")
	h.db.Model(&models.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true)
	return c.JSON(fiber.Map{"success": true})
}

// MarkAllAsRead — POST /api/notifications/read-all
func (h *NotificationsHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	result := h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = false", userID).Update("is_read", true)
	return c.JSON(fiber.Map{"success": true, "updatedCount": result.RowsAffected})
}

// GetUnreadCount — GET /api/notifications/count
func (h *NotificationsHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var count int64
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = false", userID).Count(&count)
	return c.JSON(fiber.Map{"success": true, "unreadCount": count})
}
