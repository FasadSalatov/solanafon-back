package handlers

import (
	"fmt"
	"time"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CallsHandler handles /api/calls/* endpoints
type CallsHandler struct {
	db *gorm.DB
}

func NewCallsHandler(db *gorm.DB) *CallsHandler {
	return &CallsHandler{db: db}
}

// CreateRoom — POST /api/calls/rooms/create
func (h *CallsHandler) CreateRoom(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var input struct {
		Type           string   `json:"type"`
		ParticipantIDs []string `json:"participantIds"`
	}
	c.BodyParser(&input)

	roomType := input.Type
	if roomType == "" {
		roomType = "direct"
	}

	roomCode := models.GenerateRoomCode()
	room := models.CallRoom{
		RoomCode: roomCode, Type: roomType, Status: "waiting", CreatedBy: userID,
	}
	h.db.Create(&room)

	now := time.Now()
	h.db.Create(&models.CallParticipant{
		RoomID: room.ID, UserID: userID, Status: "joined", JoinedAt: &now,
	})

	h.db.Preload("Participants").Preload("Participants.User").First(&room, room.ID)

	return c.Status(201).JSON(fiber.Map{
		"data": fiber.Map{
			"room":     formatRoom(room),
			"deepLink": fmt.Sprintf("solafon://call/%s", roomCode),
			"roomCode": roomCode,
		},
	})
}

// GetRoomByCode — GET /api/calls/rooms/code/:roomCode
func (h *CallsHandler) GetRoomByCode(c *fiber.Ctx) error {
	code := c.Params("roomCode")
	var room models.CallRoom
	if err := h.db.Where("room_code = ?", code).Preload("Participants").Preload("Participants.User").
		First(&room).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "Room not found"}})
	}
	return c.JSON(fiber.Map{"data": fiber.Map{"room": formatRoom(room)}})
}

// JoinRoom — POST /api/calls/rooms/:roomId/join
func (h *CallsHandler) JoinRoom(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	roomID := c.Params("roomId")

	var room models.CallRoom
	if err := h.db.First(&room, roomID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "Room not found"}})
	}

	now := time.Now()
	h.db.Create(&models.CallParticipant{
		RoomID: room.ID, UserID: userID, Status: "joined", JoinedAt: &now,
	})

	if room.Status == "waiting" {
		room.Status = "active"
		room.StartedAt = &now
		h.db.Save(&room)
	}

	h.db.Preload("Participants").Preload("Participants.User").First(&room, room.ID)
	return c.JSON(fiber.Map{"data": fiber.Map{"room": formatRoom(room)}})
}

// EndCall — POST /api/calls/rooms/:roomId/end
func (h *CallsHandler) EndCall(c *fiber.Ctx) error {
	roomID := c.Params("roomId")
	var room models.CallRoom
	if err := h.db.First(&room, roomID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "Room not found"}})
	}

	now := time.Now()
	room.Status = "ended"
	room.EndedAt = &now
	if room.StartedAt != nil {
		room.Duration = int(now.Sub(*room.StartedAt).Seconds())
	}
	h.db.Save(&room)

	h.db.Model(&models.CallParticipant{}).Where("room_id = ?", room.ID).Update("status", "left")

	return c.JSON(fiber.Map{"success": true})
}

// UpdateParticipantStatus — PATCH /api/calls/rooms/:roomId/status
func (h *CallsHandler) UpdateParticipantStatus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	roomID := c.Params("roomId")

	var input struct {
		IsMuted   *bool `json:"isMuted"`
		IsVideoOn *bool `json:"isVideoOn"`
		IsAudioOn *bool `json:"isAudioOn"`
	}
	c.BodyParser(&input)

	updates := map[string]interface{}{}
	if input.IsMuted != nil {
		updates["is_muted"] = *input.IsMuted
	}
	if input.IsVideoOn != nil {
		updates["is_video_on"] = *input.IsVideoOn
	}
	if input.IsAudioOn != nil {
		updates["is_audio_on"] = *input.IsAudioOn
	}

	h.db.Model(&models.CallParticipant{}).Where("room_id = ? AND user_id = ?", roomID, userID).Updates(updates)
	return c.JSON(fiber.Map{"success": true})
}

func formatRoom(room models.CallRoom) fiber.Map {
	participants := make([]fiber.Map, len(room.Participants))
	for i, p := range room.Participants {
		var userData fiber.Map
		if p.User.ID > 0 {
			userData = fiber.Map{
				"id": fmt.Sprintf("user_%d", p.User.ID),
				"displayName": p.User.GetDisplayName(), "avatarUrl": p.User.GetAvatarURL(),
			}
		}
		participants[i] = fiber.Map{
			"id": fmt.Sprintf("part_%d", p.ID), "roomId": fmt.Sprintf("room_%d", p.RoomID),
			"userId": fmt.Sprintf("user_%d", p.UserID), "status": p.Status,
			"joinedAt": p.JoinedAt, "isMuted": p.IsMuted, "isVideoOn": p.IsVideoOn,
			"isAudioOn": p.IsAudioOn, "user": userData,
		}
	}
	return fiber.Map{
		"id": fmt.Sprintf("room_%d", room.ID), "roomCode": room.RoomCode,
		"type": room.Type, "status": room.Status,
		"createdBy": fmt.Sprintf("user_%d", room.CreatedBy),
		"startedAt": room.StartedAt, "endedAt": room.EndedAt,
		"duration": room.Duration, "createdAt": room.CreatedAt,
		"participants": participants,
	}
}
