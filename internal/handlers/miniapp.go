package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MiniAppHandler struct {
	db *gorm.DB
}

func NewMiniAppHandler(db *gorm.DB) *MiniAppHandler {
	return &MiniAppHandler{db: db}
}

// GetAll - get all approved apps (with optional category filter)
func (h *MiniAppHandler) GetAll(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)
	categorySlug := c.Query("category") // Optional filter

	var user models.User
	h.db.First(&user, userID)

	query := h.db.Preload("Category").Where("moderation_status = ?", models.ModerationApproved)

	// Filter by category if provided
	if categorySlug != "" && categorySlug != "all" {
		var category models.Category
		if err := h.db.Where("slug = ?", categorySlug).First(&category).Error; err == nil {
			query = query.Where("category_id = ?", category.ID)
		}
	}

	// If user doesn't have secret access, exclude secret apps
	if !user.HasSecretAccess {
		query = query.Where("is_secret = ?", false)
	}

	var miniApps []models.MiniApp
	if err := query.Order("users_count DESC, created_at DESC").Find(&miniApps).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch mini apps",
		})
	}

	return c.JSON(fiber.Map{
		"apps":  miniApps,
		"total": len(miniApps),
	})
}

// Search - search apps by title, subtitle or description
func (h *MiniAppHandler) Search(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)
	searchQuery := strings.TrimSpace(c.Query("q"))

	if searchQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	var user models.User
	h.db.First(&user, userID)

	searchPattern := "%" + strings.ToLower(searchQuery) + "%"

	query := h.db.Preload("Category").
		Where("moderation_status = ?", models.ModerationApproved).
		Where("LOWER(title) LIKE ? OR LOWER(subtitle) LIKE ? OR LOWER(description) LIKE ?",
			searchPattern, searchPattern, searchPattern)

	if !user.HasSecretAccess {
		query = query.Where("is_secret = ?", false)
	}

	var miniApps []models.MiniApp
	if err := query.Order("users_count DESC").Limit(20).Find(&miniApps).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search apps",
		})
	}

	return c.JSON(fiber.Map{
		"apps":  miniApps,
		"total": len(miniApps),
		"query": searchQuery,
	})
}

// GetByID - get app details
func (h *MiniAppHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("userID").(uint)

	var miniApp models.MiniApp
	if err := h.db.Preload("Category").First(&miniApp, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Mini app not found",
		})
	}

	// Check if app is secret and user has access
	if miniApp.IsSecret {
		var user models.User
		h.db.First(&user, userID)

		if !user.HasSecretAccess {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Secret access required",
			})
		}
	}

	// Track app usage
	h.trackAppUsage(userID, miniApp.ID)

	return c.JSON(miniApp)
}

// GetCategories - get all categories
func (h *MiniAppHandler) GetCategories(c *fiber.Ctx) error {
	var categories []models.Category
	if err := h.db.Order("\"order\" ASC, name ASC").Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	return c.JSON(fiber.Map{
		"categories": categories,
		"total":      len(categories),
	})
}

// GetByCategory - get apps by category slug
func (h *MiniAppHandler) GetByCategory(c *fiber.Ctx) error {
	categorySlug := c.Params("slug")
	userID, _ := c.Locals("userID").(uint)

	var user models.User
	h.db.First(&user, userID)

	var category models.Category
	if err := h.db.Where("slug = ?", categorySlug).First(&category).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	query := h.db.Where("category_id = ? AND moderation_status = ?", category.ID, models.ModerationApproved)

	if !user.HasSecretAccess {
		query = query.Where("is_secret = ?", false)
	}

	var miniApps []models.MiniApp
	if err := query.Order("users_count DESC").Find(&miniApps).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch mini apps",
		})
	}

	return c.JSON(fiber.Map{
		"category": category,
		"apps":     miniApps,
		"total":    len(miniApps),
	})
}

// CreateAppInput - input for creating a new app
type CreateAppInput struct {
	Icon           string `json:"icon"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	CategoryID     uint   `json:"categoryId"`
	URL            string `json:"url,omitempty"`
	BotUsername    string `json:"botUsername,omitempty"`
	WelcomeMessage string `json:"welcomeMessage,omitempty"`
	WebhookURL     string `json:"webhookUrl,omitempty"`
}

// CreateApp - create a new mini app (user's app)
func (h *MiniAppHandler) CreateApp(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input CreateAppInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if input.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}
	if len(input.Title) > 50 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title must be 50 characters or less",
		})
	}
	if input.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Description is required",
		})
	}
	if input.CategoryID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category is required",
		})
	}
	if input.Icon == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Icon is required",
		})
	}

	// Verify category exists
	var category models.Category
	if err := h.db.First(&category, input.CategoryID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category",
		})
	}

	// Generate unique API token for the app
	apiToken := models.GenerateAPIToken()

	// Validate bot username uniqueness if provided
	if input.BotUsername != "" {
		var existingApp models.MiniApp
		if err := h.db.Where("bot_username = ?", input.BotUsername).First(&existingApp).Error; err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Bot username already taken",
			})
		}
	}

	app := models.MiniApp{
		Title:            input.Title,
		Subtitle:         input.Description,
		Description:      input.Description,
		Icon:             input.Icon,
		CategoryID:       input.CategoryID,
		URL:              input.URL,
		CreatorID:        userID,
		IsVerified:       false,
		IsSecret:         false,
		UsersCount:       0,
		ModerationStatus: models.ModerationPending,
		APIToken:         apiToken,
		BotUsername:      input.BotUsername,
		WelcomeMessage:   input.WelcomeMessage,
		WebhookURL:       input.WebhookURL,
	}

	if err := h.db.Create(&app).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create app",
		})
	}

	// Create default /start command if welcome message provided
	if input.WelcomeMessage != "" {
		startCmd := models.BotCommand{
			AppID:       app.ID,
			Command:     "/start",
			Description: "Start the bot",
			Response:    input.WelcomeMessage,
			IsEnabled:   true,
		}
		h.db.Create(&startCmd)
	}

	// Load category for response
	h.db.Preload("Category").First(&app, app.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "App created successfully. It will be reviewed by moderators within 24 hours.",
		"app":      app,
		"apiToken": apiToken,
	})
}

// GetMyApps - get apps created by current user
func (h *MiniAppHandler) GetMyApps(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var apps []models.MiniApp
	if err := h.db.Preload("Category").
		Where("creator_id = ?", userID).
		Order("created_at DESC").
		Find(&apps).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch apps",
		})
	}

	return c.JSON(fiber.Map{
		"apps":  apps,
		"total": len(apps),
	})
}

// UpdateAppInput - input for updating app
type UpdateAppInput struct {
	Icon           string `json:"icon,omitempty"`
	Title          string `json:"title,omitempty"`
	Description    string `json:"description,omitempty"`
	CategoryID     uint   `json:"categoryId,omitempty"`
	URL            string `json:"url,omitempty"`
	BotUsername    string `json:"botUsername,omitempty"`
	WelcomeMessage string `json:"welcomeMessage,omitempty"`
	WebhookURL     string `json:"webhookUrl,omitempty"`
}

// UpdateApp - update user's own app
func (h *MiniAppHandler) UpdateApp(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("id")

	var app models.MiniApp
	if err := h.db.First(&app, appID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "App not found",
		})
	}

	// Check ownership
	if app.CreatorID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only update your own apps",
		})
	}

	var input UpdateAppInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	contentChanged := false
	if input.Title != "" {
		if len(input.Title) > 50 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Title must be 50 characters or less",
			})
		}
		app.Title = input.Title
		contentChanged = true
	}
	if input.Description != "" {
		app.Description = input.Description
		app.Subtitle = input.Description
		contentChanged = true
	}
	if input.Icon != "" {
		app.Icon = input.Icon
		contentChanged = true
	}
	if input.CategoryID != 0 {
		app.CategoryID = input.CategoryID
		contentChanged = true
	}
	if input.URL != "" {
		app.URL = input.URL
	}

	// Bot settings (don't require re-moderation)
	if input.BotUsername != "" && input.BotUsername != app.BotUsername {
		// Check uniqueness
		var existingApp models.MiniApp
		if err := h.db.Where("bot_username = ? AND id != ?", input.BotUsername, app.ID).First(&existingApp).Error; err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Bot username already taken",
			})
		}
		app.BotUsername = input.BotUsername
	}
	if input.WelcomeMessage != "" {
		app.WelcomeMessage = input.WelcomeMessage
	}
	if input.WebhookURL != "" {
		app.WebhookURL = input.WebhookURL
	}

	// Reset moderation status only if content changed
	if contentChanged {
		app.ModerationStatus = models.ModerationPending
	}

	if err := h.db.Save(&app).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update app",
		})
	}

	h.db.Preload("Category").First(&app, app.ID)

	return c.JSON(fiber.Map{
		"message": "App updated. It will be re-reviewed by moderators.",
		"app":     app,
	})
}

// DeleteApp - delete user's own app
func (h *MiniAppHandler) DeleteApp(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("id")

	var app models.MiniApp
	if err := h.db.First(&app, appID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "App not found",
		})
	}

	// Check ownership
	if app.CreatorID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only delete your own apps",
		})
	}

	if err := h.db.Delete(&app).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete app",
		})
	}

	return c.JSON(fiber.Map{
		"message": "App deleted successfully",
	})
}

// Helper: track app usage
func (h *MiniAppHandler) trackAppUsage(userID, appID uint) {
	var appUser models.AppUser
	result := h.db.Where("user_id = ? AND app_id = ?", userID, appID).First(&appUser)

	if result.Error == gorm.ErrRecordNotFound {
		// New user for this app
		appUser = models.AppUser{
			UserID:    userID,
			AppID:     appID,
			LastUsed:  time.Now(),
			CreatedAt: time.Now(),
		}
		h.db.Create(&appUser)

		// Increment users count
		h.db.Model(&models.MiniApp{}).Where("id = ?", appID).
			UpdateColumn("users_count", gorm.Expr("users_count + 1"))
	} else {
		// Update last used
		h.db.Model(&appUser).Update("last_used", time.Now())
	}
}

// GetAppMessages - get chat messages for an app
func (h *MiniAppHandler) GetAppMessages(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("id")

	var messages []models.AppMessage
	if err := h.db.Where("app_id = ? AND user_id = ?", appID, userID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch messages",
		})
	}

	// Mark as read
	h.db.Model(&models.AppMessage{}).
		Where("app_id = ? AND user_id = ? AND is_read = ?", appID, userID, false).
		Update("is_read", true)

	return c.JSON(fiber.Map{
		"messages": messages,
		"total":    len(messages),
	})
}

// SendAppMessageInput - input for sending message
type SendAppMessageInput struct {
	Content string `json:"content"`
}

// SendAppMessage - send message to app (triggers bot response or webhook)
func (h *MiniAppHandler) SendAppMessage(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("id")

	var input SendAppMessageInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if strings.TrimSpace(input.Content) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message content is required",
		})
	}

	// Verify app exists
	var app models.MiniApp
	if err := h.db.First(&app, appID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "App not found",
		})
	}

	// Get user info
	var user models.User
	h.db.First(&user, userID)

	// Create user message
	userMsg := models.AppMessage{
		AppID:       app.ID,
		UserID:      userID,
		Content:     input.Content,
		IsFromBot:   false,
		IsRead:      true,
		MessageType: "text",
		CreatedAt:   time.Now(),
	}
	h.db.Create(&userMsg)

	// Track app usage
	h.trackAppUsage(userID, app.ID)

	// Try to get bot response
	botResponse := h.getBotResponse(app, user, input.Content, userMsg.ID)

	var botMsg *models.AppMessage
	if botResponse != "" {
		botMsg = &models.AppMessage{
			AppID:       app.ID,
			UserID:      userID,
			Content:     botResponse,
			IsFromBot:   true,
			IsRead:      false,
			MessageType: "text",
			CreatedAt:   time.Now(),
		}
		h.db.Create(botMsg)
	}

	response := fiber.Map{
		"userMessage": userMsg,
	}
	if botMsg != nil {
		response["botMessage"] = botMsg
	}

	return c.JSON(response)
}

// getBotResponse - get response from bot command or webhook
func (h *MiniAppHandler) getBotResponse(app models.MiniApp, user models.User, message string, messageID uint) string {
	// Check for predefined command
	if strings.HasPrefix(message, "/") {
		var cmd models.BotCommand
		if err := h.db.Where("app_id = ? AND command = ? AND is_enabled = ?", app.ID, message, true).First(&cmd).Error; err == nil {
			return cmd.Response
		}
	}

	// Check for /start with welcome message
	if strings.ToLower(message) == "/start" && app.WelcomeMessage != "" {
		return app.WelcomeMessage
	}

	// If webhook is configured, trigger it asynchronously
	if app.WebhookURL != "" {
		go h.triggerWebhook(app, user, message, messageID)
		return "" // Response will come later via Bot API
	}

	// Default response if no webhook
	if strings.ToLower(message) == "/start" {
		return "Привет! 👋\n\nДобро пожаловать в " + app.Title + "!\n\n" + app.Description
	}

	return ""
}

// triggerWebhook - send message to app's webhook
func (h *MiniAppHandler) triggerWebhook(app models.MiniApp, user models.User, message string, messageID uint) {
	startTime := time.Now()

	payload := map[string]interface{}{
		"update_id":  messageID,
		"message_id": messageID,
		"from": map[string]interface{}{
			"id":       user.ID,
			"email":    user.Email,
			"name":     user.Name,
			"language": user.Language,
		},
		"chat": map[string]interface{}{
			"id":   user.ID,
			"type": "private",
		},
		"date": time.Now().Unix(),
		"text": message,
	}

	payloadBytes, _ := json.Marshal(payload)

	webhookLog := models.WebhookLog{
		AppID:      app.ID,
		URL:        app.WebhookURL,
		Method:     "POST",
		StatusCode: 0,
		Request:    string(payloadBytes),
		Response:   "",
		Duration:   0,
		CreatedAt:  time.Now(),
	}

	// Make actual HTTP request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("POST", app.WebhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		webhookLog.Response = "Error creating request: " + err.Error()
		webhookLog.Duration = int(time.Since(startTime).Milliseconds())
		h.db.Create(&webhookLog)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-ID", fmt.Sprintf("%d", app.ID))
	req.Header.Set("X-App-Token", app.APIToken)

	resp, err := client.Do(req)
	if err != nil {
		webhookLog.Response = "Error sending request: " + err.Error()
		webhookLog.Duration = int(time.Since(startTime).Milliseconds())
		h.db.Create(&webhookLog)
		return
	}
	defer resp.Body.Close()

	// Read response body
	var respBody bytes.Buffer
	respBody.ReadFrom(resp.Body)

	webhookLog.StatusCode = resp.StatusCode
	webhookLog.Response = respBody.String()
	webhookLog.Duration = int(time.Since(startTime).Milliseconds())

	h.db.Create(&webhookLog)
}

// RegenerateAPIToken - regenerate API token for an app
func (h *MiniAppHandler) RegenerateAPIToken(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("id")

	var app models.MiniApp
	if err := h.db.First(&app, appID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "App not found",
		})
	}

	if app.CreatorID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only manage your own apps",
		})
	}

	newToken := models.GenerateAPIToken()
	app.APIToken = newToken

	if err := h.db.Save(&app).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to regenerate token",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "API token regenerated successfully",
		"apiToken": newToken,
	})
}

// GetAppSettings - get app settings including API token (for owner only)
func (h *MiniAppHandler) GetAppSettings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("id")

	var app models.MiniApp
	if err := h.db.Preload("Category").First(&app, appID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "App not found",
		})
	}

	if app.CreatorID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only view settings for your own apps",
		})
	}

	// Get bot commands
	var commands []models.BotCommand
	h.db.Where("app_id = ?", app.ID).Order("command ASC").Find(&commands)

	// Get recent webhook logs
	var webhookLogs []models.WebhookLog
	h.db.Where("app_id = ?", app.ID).Order("created_at DESC").Limit(10).Find(&webhookLogs)

	return c.JSON(fiber.Map{
		"app": fiber.Map{
			"id":               app.ID,
			"title":            app.Title,
			"description":      app.Description,
			"icon":             app.Icon,
			"categoryId":       app.CategoryID,
			"category":         app.Category,
			"url":              app.URL,
			"botUsername":      app.BotUsername,
			"welcomeMessage":   app.WelcomeMessage,
			"webhookUrl":       app.WebhookURL,
			"apiToken":         app.APIToken,
			"moderationStatus": app.ModerationStatus,
			"usersCount":       app.UsersCount,
			"createdAt":        app.CreatedAt,
		},
		"commands":    commands,
		"webhookLogs": webhookLogs,
	})
}

// BotCommandInput - input for bot commands
type BotCommandInput struct {
	Command     string `json:"command"`
	Description string `json:"description"`
	Response    string `json:"response"`
	IsEnabled   bool   `json:"isEnabled"`
}

// AddBotCommand - add a bot command
func (h *MiniAppHandler) AddBotCommand(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("id")

	var app models.MiniApp
	if err := h.db.First(&app, appID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "App not found",
		})
	}

	if app.CreatorID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only manage your own apps",
		})
	}

	var input BotCommandInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.Command == "" || !strings.HasPrefix(input.Command, "/") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Command must start with /",
		})
	}

	// Check if command already exists
	var existingCmd models.BotCommand
	if err := h.db.Where("app_id = ? AND command = ?", app.ID, input.Command).First(&existingCmd).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Command already exists",
		})
	}

	cmd := models.BotCommand{
		AppID:       app.ID,
		Command:     input.Command,
		Description: input.Description,
		Response:    input.Response,
		IsEnabled:   true,
	}

	if err := h.db.Create(&cmd).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create command",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Command created successfully",
		"command": cmd,
	})
}

// UpdateBotCommand - update a bot command
func (h *MiniAppHandler) UpdateBotCommand(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("id")
	cmdID := c.Params("cmdId")

	var app models.MiniApp
	if err := h.db.First(&app, appID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "App not found",
		})
	}

	if app.CreatorID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only manage your own apps",
		})
	}

	var cmd models.BotCommand
	if err := h.db.Where("id = ? AND app_id = ?", cmdID, app.ID).First(&cmd).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Command not found",
		})
	}

	var input BotCommandInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.Description != "" {
		cmd.Description = input.Description
	}
	if input.Response != "" {
		cmd.Response = input.Response
	}
	cmd.IsEnabled = input.IsEnabled

	if err := h.db.Save(&cmd).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update command",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Command updated successfully",
		"command": cmd,
	})
}

// DeleteBotCommand - delete a bot command
func (h *MiniAppHandler) DeleteBotCommand(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID := c.Params("id")
	cmdID := c.Params("cmdId")

	var app models.MiniApp
	if err := h.db.First(&app, appID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "App not found",
		})
	}

	if app.CreatorID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only manage your own apps",
		})
	}

	if err := h.db.Where("id = ? AND app_id = ?", cmdID, app.ID).Delete(&models.BotCommand{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete command",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Command deleted successfully",
	})
}

// GetBotCommands - get all commands for an app
func (h *MiniAppHandler) GetBotCommands(c *fiber.Ctx) error {
	appID := c.Params("id")

	var commands []models.BotCommand
	if err := h.db.Where("app_id = ? AND is_enabled = ?", appID, true).
		Order("command ASC").Find(&commands).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch commands",
		})
	}

	return c.JSON(fiber.Map{
		"commands": commands,
		"total":    len(commands),
	})
}
