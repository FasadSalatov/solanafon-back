package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CrashHandler handles /api/crash/* endpoints
type CrashHandler struct {
	db *gorm.DB
}

func NewCrashHandler(db *gorm.DB) *CrashHandler {
	return &CrashHandler{db: db}
}

// ReportCrash — POST /api/crash/mobile
func (h *CrashHandler) ReportCrash(c *fiber.Ctx) error {
	var input struct {
		UserID    string          `json:"userId"`
		SessionID string          `json:"sessionId"`
		ErrorType string          `json:"errorType"`
		Message   string          `json:"message"`
		Stack     string          `json:"stack"`
		UserAgent string          `json:"userAgent"`
		Source    string          `json:"source"`
		Metadata  json.RawMessage `json:"metadata"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "Invalid body"}})
	}

	if input.ErrorType == "" {
		input.ErrorType = "error"
	}
	if input.Source == "" {
		input.Source = "android"
	}

	report := models.CrashReport{
		UserID: input.UserID, SessionID: input.SessionID,
		ErrorType: input.ErrorType, Message: input.Message,
		Stack: input.Stack, UserAgent: input.UserAgent,
		Source: input.Source, Metadata: string(input.Metadata),
	}
	h.db.Create(&report)

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"id": fmt.Sprintf("crash_%d", report.ID), "appId": "solafon-android",
			"errorType": report.ErrorType, "message": report.Message, "createdAt": report.CreatedAt,
		},
	})
}
