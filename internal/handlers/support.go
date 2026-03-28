package handlers

import (
	"fmt"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SupportHandler handles /api/support/*, /api/legal/*, /api/i18n/*
type SupportHandler struct {
	db *gorm.DB
}

func NewSupportHandler(db *gorm.DB) *SupportHandler {
	return &SupportHandler{db: db}
}

// GetFAQ — GET /api/support/faq
func (h *SupportHandler) GetFAQ(c *fiber.Ctx) error {
	lang := c.Query("language", "en")
	var faqs []models.FAQ
	h.db.Where("language = ?", lang).Order("`order` ASC").Find(&faqs)

	if len(faqs) == 0 {
		// Seed default FAQs
		defaults := []models.FAQ{
			{Question: "How to create a wallet?", Answer: "Go to the Wallet tab and tap 'Create Wallet'. A 12-word seed phrase will be generated. Write it down and store it safely.", Category: "wallet", Language: "en", Order: 1},
			{Question: "How to send tokens?", Answer: "Go to Wallet → Send. Enter the recipient address, amount, and confirm the transaction.", Category: "wallet", Language: "en", Order: 2},
			{Question: "What are Mana Points?", Answer: "Mana Points (MP) is the internal currency used for premium features, Secret Login, and other services on the platform.", Category: "general", Language: "en", Order: 3},
			{Question: "How to create an app?", Answer: "Go to Developer section or use Dev Studio to create a new mini-app. Your app will go through moderation before being published.", Category: "developer", Language: "en", Order: 4},
			{Question: "Как создать кошелёк?", Answer: "Перейдите во вкладку Кошелёк и нажмите 'Создать кошелёк'. Будет сгенерирована фраза из 12 слов. Запишите её и храните в безопасном месте.", Category: "wallet", Language: "ru", Order: 1},
			{Question: "Что такое Mana Points?", Answer: "Mana Points (MP) — внутренняя валюта для премиум-функций, Secret Login и других сервисов на платформе.", Category: "general", Language: "ru", Order: 2},
		}
		for i := range defaults {
			h.db.Create(&defaults[i])
		}
		h.db.Where("language = ?", lang).Order("`order` ASC").Find(&faqs)
	}

	result := make([]fiber.Map, len(faqs))
	for i, f := range faqs {
		result[i] = fiber.Map{
			"id": fmt.Sprintf("faq_%d", f.ID), "question": f.Question,
			"answer": f.Answer, "category": f.Category,
		}
	}
	return c.JSON(fiber.Map{"faq": result})
}

// CreateTicket — POST /api/support/tickets
func (h *SupportHandler) CreateTicket(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var input struct {
		Subject  string `json:"subject"`
		Message  string `json:"message"`
		Category string `json:"category"`
	}
	if err := c.BodyParser(&input); err != nil || input.Subject == "" || input.Message == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "subject and message are required"}})
	}

	cat := input.Category
	if cat == "" {
		cat = "general"
	}
	ticket := models.SupportTicket{UserID: userID, Subject: input.Subject, Message: input.Message, Category: cat}
	h.db.Create(&ticket)

	return c.Status(201).JSON(fiber.Map{
		"ticketId": fmt.Sprintf("ticket_%d", ticket.ID), "status": "open",
	})
}

// GetTerms — GET /api/legal/terms
func (h *SupportHandler) GetTerms(c *fiber.Ctx) error {
	return h.getLegalDoc(c, "terms")
}

// GetPrivacy — GET /api/legal/privacy
func (h *SupportHandler) GetPrivacy(c *fiber.Ctx) error {
	return h.getLegalDoc(c, "privacy")
}

func (h *SupportHandler) getLegalDoc(c *fiber.Ctx, docType string) error {
	var doc models.LegalDocument
	if err := h.db.Where("type = ?", docType).First(&doc).Error; err != nil {
		// Create default
		content := "# Terms of Service\n\nLast updated: 2026-03-01\n\nBy using Solafon, you agree to these terms..."
		if docType == "privacy" {
			content = "# Privacy Policy\n\nLast updated: 2026-03-01\n\nYour privacy is important to us..."
		}
		doc = models.LegalDocument{Type: docType, Content: content, Version: "1.0", LastUpdated: "2026-03-01"}
		h.db.Create(&doc)
	}
	return c.JSON(fiber.Map{"content": doc.Content, "version": doc.Version, "lastUpdated": doc.LastUpdated})
}

// GetLanguages — GET /api/i18n/languages
func (h *SupportHandler) GetLanguages(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"languages": []fiber.Map{
			{"code": "en", "name": "English", "flag": "🇺🇸"},
			{"code": "ru", "name": "Русский", "flag": "🇷🇺"},
			{"code": "zh", "name": "中文", "flag": "🇨🇳"},
			{"code": "es", "name": "Español", "flag": "🇪🇸"},
		},
		"defaultLanguage": "en",
	})
}
