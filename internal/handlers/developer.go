package handlers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fasad/solanafon-back/internal/config"
	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DeveloperHandler handles /api/developer/* endpoints
type DeveloperHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewDeveloperHandler(db *gorm.DB, cfg *config.Config) *DeveloperHandler {
	return &DeveloperHandler{db: db, cfg: cfg}
}

// ListMyApps — GET /api/developer/apps
func (h *DeveloperHandler) ListMyApps(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var apps []models.MiniApp
	h.db.Where("creator_id = ?", userID).Preload("Category").Find(&apps)

	var totalUsers int
	var approved, pending, rejected int64
	for _, a := range apps {
		totalUsers += a.UsersCount
		switch a.ModerationStatus {
		case models.ModerationApproved:
			approved++
		case models.ModerationPending:
			pending++
		case models.ModerationRejected:
			rejected++
		}
	}

	result := make([]fiber.Map, len(apps))
	for i, a := range apps {
		result[i] = fiber.Map{
			"id": fmt.Sprintf("app_%d", a.ID), "name": a.Title,
			"description": a.Description, "icon": a.Icon, "iconUrl": a.IconURL,
			"category": a.Category.Slug, "url": a.URL,
			"status": string(a.ModerationStatus), "isVerified": a.IsVerified,
			"users": a.FormatUsersCount(), "usersCount": a.UsersCount,
			"rating": a.Rating, "createdAt": a.CreatedAt, "updatedAt": a.UpdatedAt,
			"moderationNote": a.ModerationNote,
		}
	}

	return c.JSON(fiber.Map{
		"apps": result,
		"stats": fiber.Map{
			"totalApps": len(apps), "approvedApps": approved,
			"pendingApps": pending, "rejectedApps": rejected, "totalUsers": totalUsers,
		},
	})
}

// CreateApp — POST /api/developer/apps
func (h *DeveloperHandler) CreateApp(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var input struct {
		Name     string `json:"name"`
		Description string `json:"description"`
		Icon     string `json:"icon"`
		IconURL  string `json:"iconUrl"`
		Category string `json:"category"`
		URL      string `json:"url"`
	}
	if err := c.BodyParser(&input); err != nil || input.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "name is required"}})
	}

	// Resolve category
	var category models.Category
	h.db.Where("slug = ? OR name = ?", input.Category, input.Category).First(&category)
	if category.ID == 0 {
		category.ID = 1 // default
	}

	apiKey := models.GenerateAPIToken()
	webhookSecret := generateWebhookSecret()

	app := models.MiniApp{
		Title: input.Name, Description: input.Description, Icon: input.Icon, IconURL: input.IconURL,
		CategoryID: category.ID, URL: input.URL, CreatorID: userID,
		APIToken: apiKey, WebhookSecret: webhookSecret,
		ModerationStatus: models.ModerationPending,
	}
	h.db.Create(&app)

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"app":     formatDevApp(app),
		"message": "App created successfully",
		"apiKey":  apiKey,
		"keyId":   fmt.Sprintf("key_%d", app.ID),
		"apiKeyHint": apiKey[:12] + "...",
	})
}

// GetApp — GET /api/developer/apps/:appId
func (h *DeveloperHandler) GetApp(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("appId")
	var app models.MiniApp
	if err := h.db.Where("id = ? AND creator_id = ?", appID, userID).Preload("Category").First(&app).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "App not found"}})
	}
	return c.JSON(formatDevApp(app))
}

// UpdateApp — PATCH /api/developer/apps/:appId
func (h *DeveloperHandler) UpdateApp(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("appId")
	var app models.MiniApp
	if err := h.db.Where("id = ? AND creator_id = ?", appID, userID).First(&app).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "App not found"}})
	}

	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Icon        *string `json:"icon"`
		IconURL     *string `json:"iconUrl"`
		Category    *string `json:"category"`
		URL         *string `json:"url"`
	}
	c.BodyParser(&input)

	needsRemoderation := false
	if input.Name != nil {
		app.Title = *input.Name
		needsRemoderation = true
	}
	if input.Description != nil {
		app.Description = *input.Description
		needsRemoderation = true
	}
	if input.Icon != nil {
		app.Icon = *input.Icon
	}
	if input.IconURL != nil {
		app.IconURL = *input.IconURL
	}
	if input.URL != nil {
		app.URL = *input.URL
	}
	if input.Category != nil {
		var cat models.Category
		if h.db.Where("slug = ? OR name = ?", *input.Category, *input.Category).First(&cat).Error == nil {
			app.CategoryID = cat.ID
		}
	}

	if needsRemoderation {
		app.ModerationStatus = models.ModerationPending
	}

	h.db.Save(&app)
	return c.JSON(fiber.Map{"success": true, "app": formatDevApp(app)})
}

// DeleteApp — DELETE /api/developer/apps/:appId
func (h *DeveloperHandler) DeleteApp(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("appId")
	h.db.Where("id = ? AND creator_id = ?", appID, userID).Delete(&models.MiniApp{})
	return c.JSON(fiber.Map{"success": true})
}

// GenerateAPIKey — POST /api/developer/apps/:appId/api-key
func (h *DeveloperHandler) GenerateAPIKey(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("appId")
	var app models.MiniApp
	if err := h.db.Where("id = ? AND creator_id = ?", appID, userID).First(&app).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "App not found"}})
	}

	newKey := models.GenerateAPIToken()
	newSecret := generateWebhookSecret()
	app.APIToken = newKey
	app.WebhookSecret = newSecret
	h.db.Save(&app)

	return c.JSON(fiber.Map{
		"success": true, "apiKey": newKey, "apiSecret": newSecret,
		"keyId": fmt.Sprintf("key_%d", app.ID), "message": "Key generated",
	})
}

// ListAPICredentials — GET /api/developer/apps/:appId/api-key
func (h *DeveloperHandler) ListAPICredentials(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("appId")
	var app models.MiniApp
	if err := h.db.Where("id = ? AND creator_id = ?", appID, userID).First(&app).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "App not found"}})
	}

	hint := ""
	if len(app.APIToken) > 12 {
		hint = app.APIToken[:12] + "..."
	}

	return c.JSON(fiber.Map{
		"credentials": []fiber.Map{{
			"id": fmt.Sprintf("key_%d", app.ID), "appId": fmt.Sprintf("app_%d", app.ID),
			"apiKeyPrefix": hint, "webhookUrl": app.WebhookURL,
			"webhookSecret": app.WebhookSecret, "isActive": true,
			"createdAt": app.CreatedAt,
		}},
	})
}

// RevokeAPIKey — DELETE /api/developer/apps/:appId/credentials/:keyId
func (h *DeveloperHandler) RevokeAPIKey(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("appId")
	var app models.MiniApp
	if err := h.db.Where("id = ? AND creator_id = ?", appID, userID).First(&app).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "App not found"}})
	}
	app.APIToken = ""
	h.db.Save(&app)
	return c.JSON(fiber.Map{"success": true})
}

// UpdateWebhook — PUT /api/developer/apps/:appId/webhook
func (h *DeveloperHandler) UpdateWebhook(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("appId")
	var app models.MiniApp
	if err := h.db.Where("id = ? AND creator_id = ?", appID, userID).First(&app).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "App not found"}})
	}

	var input struct {
		WebhookURL string   `json:"webhookUrl"`
		Events     []string `json:"events"`
	}
	c.BodyParser(&input)

	app.WebhookURL = input.WebhookURL
	if app.WebhookSecret == "" {
		app.WebhookSecret = generateWebhookSecret()
	}
	h.db.Save(&app)

	return c.JSON(fiber.Map{"success": true, "webhookSecret": app.WebhookSecret})
}

// GetWelcomeMessage — GET /api/developer/apps/:appId/welcome-message
func (h *DeveloperHandler) GetWelcomeMessage(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("appId")
	var app models.MiniApp
	if err := h.db.Where("id = ? AND creator_id = ?", appID, userID).First(&app).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "App not found"}})
	}

	return c.JSON(fiber.Map{
		"welcomeMessage": fiber.Map{
			"id": fmt.Sprintf("wm_%d", app.ID), "appId": fmt.Sprintf("app_%d", app.ID),
			"content": fiber.Map{"type": "text", "text": app.WelcomeMessage},
			"isActive": app.WelcomeMessage != "", "createdAt": app.CreatedAt,
		},
	})
}

// UpdateWelcomeMessage — PUT /api/developer/apps/:appId/welcome-message
func (h *DeveloperHandler) UpdateWelcomeMessage(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("appId")
	var app models.MiniApp
	if err := h.db.Where("id = ? AND creator_id = ?", appID, userID).First(&app).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "App not found"}})
	}

	var input struct {
		Content struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		IsActive bool `json:"isActive"`
	}
	c.BodyParser(&input)

	if input.IsActive {
		app.WelcomeMessage = input.Content.Text
	} else {
		app.WelcomeMessage = ""
	}
	h.db.Save(&app)

	return c.JSON(fiber.Map{"success": true})
}

// Upload — POST /api/developer/upload
func (h *DeveloperHandler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "file is required"}})
	}

	// Validate size (5MB max)
	if file.Size > 5*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "FILE_TOO_LARGE", "message": "Max file size is 5MB"}})
	}

	// Validate type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExts := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".gif": true}
	if !validExts[ext] {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "INVALID_FORMAT", "message": "Allowed: PNG, JPG, WEBP, GIF"}})
	}

	// Generate unique filename
	b := make([]byte, 16)
	rand.Read(b)
	filename := hex.EncodeToString(b) + ext

	// Ensure upload dir exists
	uploadDir := h.cfg.UploadDir
	os.MkdirAll(uploadDir, 0755)

	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveFile(file, dst); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fiber.Map{"code": "UPLOAD_FAILED", "message": "Failed to save file"}})
	}

	url := fmt.Sprintf("%s/uploads/%s", h.cfg.BaseURL, filename)
	return c.JSON(fiber.Map{
		"url": url, "filename": filename,
		"size": file.Size, "file_type": file.Header.Get("Content-Type"),
	})
}

// Apps marketplace — enhanced endpoints

// ListApps — GET /api/apps (enhanced with pagination/sorting)
func (h *DeveloperHandler) ListApps(c *fiber.Ctx) error {
	category := c.Query("category")
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	sortBy := c.Query("sortBy", "popular")
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	query := h.db.Model(&models.MiniApp{}).Where("moderation_status = ?", "approved")

	if category != "" && category != "all" {
		var cat models.Category
		if h.db.Where("slug = ?", category).First(&cat).Error == nil {
			query = query.Where("category_id = ?", cat.ID)
		}
	}
	if search != "" {
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?",
			"%"+strings.ToLower(search)+"%", "%"+strings.ToLower(search)+"%")
	}

	var total int64
	query.Count(&total)

	switch sortBy {
	case "new":
		query = query.Order("created_at DESC")
	case "trending":
		query = query.Where("is_trending = true").Order("users_count DESC")
	default: // popular
		query = query.Order("users_count DESC")
	}

	var apps []models.MiniApp
	query.Preload("Category").Preload("Creator").Offset(offset).Limit(limit).Find(&apps)

	result := make([]fiber.Map, len(apps))
	for i, a := range apps {
		var dev fiber.Map
		if a.Creator != nil {
			dev = fiber.Map{
				"id": fmt.Sprintf("dev_%d", a.Creator.ID),
				"name": a.Creator.GetDisplayName(), "isVerified": false,
			}
		}
		result[i] = fiber.Map{
			"id": fmt.Sprintf("app_%d", a.ID), "name": a.Title,
			"description": a.Description, "icon": a.Icon, "iconUrl": a.IconURL,
			"category": a.Category.Slug, "url": a.URL,
			"users": a.FormatUsersCount(), "usersCount": a.UsersCount,
			"isVerified": a.IsVerified, "isTrending": a.IsTrending,
			"rating": a.Rating, "createdAt": a.CreatedAt, "developer": dev,
		}
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return c.JSON(fiber.Map{
		"apps": result,
		"pagination": fiber.Map{
			"currentPage": page, "totalPages": totalPages,
			"totalItems": total, "hasMore": page < totalPages,
		},
	})
}

// GetCategories — GET /api/apps/categories
func (h *DeveloperHandler) GetCategories(c *fiber.Ctx) error {
	var categories []models.Category
	h.db.Order("`order` ASC").Find(&categories)

	result := make([]fiber.Map, len(categories))
	for i, cat := range categories {
		var count int64
		h.db.Model(&models.MiniApp{}).Where("category_id = ? AND moderation_status = ?", cat.ID, "approved").Count(&count)
		result[i] = fiber.Map{
			"id": cat.Slug, "name": cat.Name, "icon": cat.Icon, "appsCount": count,
		}
	}
	return c.JSON(fiber.Map{"categories": result})
}

// GetAppDetail — GET /api/apps/:appId
func (h *DeveloperHandler) GetAppDetail(c *fiber.Ctx) error {
	appID := c.Params("appId")
	var app models.MiniApp
	if err := h.db.Preload("Category").Preload("Creator").First(&app, appID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "App not found"}})
	}

	var dev fiber.Map
	if app.Creator != nil {
		dev = fiber.Map{
			"id": fmt.Sprintf("dev_%d", app.Creator.ID),
			"name": app.Creator.GetDisplayName(), "isVerified": false,
		}
	}

	return c.JSON(fiber.Map{
		"id": fmt.Sprintf("app_%d", app.ID), "name": app.Title,
		"description": app.Description, "longDescription": app.LongDescription,
		"icon": app.Icon, "iconUrl": app.IconURL,
		"category": app.Category.Slug, "url": app.URL,
		"users": app.FormatUsersCount(), "usersCount": app.UsersCount,
		"isVerified": app.IsVerified, "isTrending": app.IsTrending,
		"rating": app.Rating, "screenshots": app.Screenshots,
		"tags": app.Tags, "reviewsCount": app.ReviewsCount,
		"createdAt": app.CreatedAt, "updatedAt": app.UpdatedAt,
		"permissions": app.Permissions, "version": app.Version,
		"developer": dev,
	})
}

// LaunchApp — POST /api/apps/:appId/launch
func (h *DeveloperHandler) LaunchApp(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID, _ := strconv.Atoi(c.Params("appId"))

	var appUser models.AppUser
	if h.db.Where("user_id = ? AND app_id = ?", userID, appID).First(&appUser).Error != nil {
		h.db.Create(&models.AppUser{UserID: userID, AppID: uint(appID)})
		h.db.Model(&models.MiniApp{}).Where("id = ?", appID).UpdateColumn("users_count", gorm.Expr("users_count + 1"))
	}

	return c.JSON(fiber.Map{"success": true, "sessionId": fmt.Sprintf("session_%d", appID)})
}

// helpers

func formatDevApp(app models.MiniApp) fiber.Map {
	return fiber.Map{
		"id": fmt.Sprintf("app_%d", app.ID), "name": app.Title,
		"description": app.Description, "icon": app.Icon, "iconUrl": app.IconURL,
		"category": app.Category.Slug, "url": app.URL,
		"status": string(app.ModerationStatus), "isVerified": app.IsVerified,
		"users": app.FormatUsersCount(), "usersCount": app.UsersCount,
		"rating": app.Rating, "createdAt": app.CreatedAt, "updatedAt": app.UpdatedAt,
		"moderationNote": app.ModerationNote,
	}
}

func generateWebhookSecret() string {
	b := make([]byte, 32)
	rand.Read(b)
	return "whsec_" + hex.EncodeToString(b)
}

// SignWebhookPayload creates HMAC-SHA256 signature for webhook payloads
func SignWebhookPayload(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// Suppress unused import
var _ = io.Discard
