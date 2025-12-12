package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DevStudioHandler - handles Dev Studio interactions for creating/managing apps
type DevStudioHandler struct {
	db *gorm.DB
}

func NewDevStudioHandler(db *gorm.DB) *DevStudioHandler {
	return &DevStudioHandler{db: db}
}


// Dev Studio conversation states
const (
	StateIdle              = "idle"
	StateAwaitingAppName   = "awaiting_app_name"
	StateAwaitingAppDesc   = "awaiting_app_desc"
	StateAwaitingAppIcon   = "awaiting_app_icon"
	StateAwaitingCategory  = "awaiting_category"
	StateAwaitingUsername  = "awaiting_username"
	StateAwaitingWelcome   = "awaiting_welcome"
	StateAwaitingWebhook   = "awaiting_webhook"
	StateSelectingApp      = "selecting_app"
	StateEditingApp        = "editing_app"
	StateAwaitingNewName   = "awaiting_new_name"
	StateAwaitingNewDesc   = "awaiting_new_desc"
	StateAwaitingCommand   = "awaiting_command"
	StateAwaitingCmdDesc   = "awaiting_cmd_desc"
	StateAwaitingCmdResp   = "awaiting_cmd_response"
	StateDeletingApp       = "deleting_app"
)

// Dev Studio commands
const (
	CmdStart       = "/start"
	CmdNewApp      = "/newapp"
	CmdMyApps      = "/myapps"
	CmdEditApp     = "/edit"
	CmdDeleteApp   = "/delete"
	CmdToken       = "/token"
	CmdCommands    = "/commands"
	CmdWebhook     = "/webhook"
	CmdHelp        = "/help"
	CmdCancel      = "/cancel"
)

// ProcessMessage - main entry point for Dev Studio messages
func (h *DevStudioHandler) ProcessMessage(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	// Find Dev Studio app
	var app models.MiniApp
	if err := h.db.Where("bot_username = ?", "devstudio").First(&app).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Dev Studio not found",
		})
	}

	var input struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	message := strings.TrimSpace(input.Content)

	// Get or create conversation state
	state := h.getState(userID)

	// Save user message
	userMsg := models.AppMessage{
		AppID:       app.ID,
		UserID:      userID,
		Content:     message,
		IsFromBot:   false,
		IsRead:      true,
		MessageType: "text",
		CreatedAt:   time.Now(),
	}
	h.db.Create(&userMsg)

	// Process message and get response
	response := h.processCommand(userID, message, state)

	// Save bot response
	botMsg := models.AppMessage{
		AppID:       app.ID,
		UserID:      userID,
		Content:     response,
		IsFromBot:   true,
		IsRead:      false,
		MessageType: "text",
		CreatedAt:   time.Now(),
	}
	h.db.Create(&botMsg)

	return c.JSON(fiber.Map{
		"userMessage": userMsg,
		"botMessage":  botMsg,
	})
}

// processCommand - process user command/message
func (h *DevStudioHandler) processCommand(userID uint, message string, state *models.ConversationState) string {
	// Handle /cancel anytime
	if message == CmdCancel {
		h.setState(userID, StateIdle, "")
		return "Действие отменено. Чем могу помочь?\n\nИспользуй /help для списка команд."
	}

	// Handle commands
	if strings.HasPrefix(message, "/") {
		return h.handleCommand(userID, message, state)
	}

	// Handle conversation states
	return h.handleState(userID, message, state)
}

// handleCommand - handle slash commands
func (h *DevStudioHandler) handleCommand(userID uint, command string, state *models.ConversationState) string {
	switch command {
	case CmdStart:
		return h.cmdStart()
	case CmdHelp:
		return h.cmdHelp()
	case CmdNewApp:
		return h.cmdNewApp(userID)
	case CmdMyApps:
		return h.cmdMyApps(userID)
	case CmdToken:
		return h.cmdToken(userID)
	case CmdEditApp:
		return h.cmdEditApp(userID)
	case CmdDeleteApp:
		return h.cmdDeleteApp(userID)
	case CmdCommands:
		return h.cmdCommands(userID)
	case CmdWebhook:
		return h.cmdWebhook(userID)
	default:
		return "Неизвестная команда. Используй /help для списка команд."
	}
}

// handleState - handle conversation state
func (h *DevStudioHandler) handleState(userID uint, message string, state *models.ConversationState) string {
	switch state.State {
	case StateAwaitingAppName:
		return h.handleAppName(userID, message)
	case StateAwaitingAppDesc:
		return h.handleAppDesc(userID, message)
	case StateAwaitingAppIcon:
		return h.handleAppIcon(userID, message)
	case StateAwaitingCategory:
		return h.handleCategory(userID, message)
	case StateAwaitingUsername:
		return h.handleUsername(userID, message)
	case StateAwaitingWelcome:
		return h.handleWelcome(userID, message)
	case StateAwaitingWebhook:
		return h.handleWebhookUrl(userID, message)
	case StateSelectingApp:
		return h.handleAppSelection(userID, message)
	case StateEditingApp:
		return h.handleEditChoice(userID, message)
	case StateAwaitingNewName:
		return h.handleNewName(userID, message)
	case StateAwaitingNewDesc:
		return h.handleNewDesc(userID, message)
	case StateAwaitingCommand:
		return h.handleNewCommand(userID, message)
	case StateAwaitingCmdDesc:
		return h.handleCmdDescription(userID, message)
	case StateAwaitingCmdResp:
		return h.handleCmdResponse(userID, message)
	case StateDeletingApp:
		return h.handleDeleteConfirm(userID, message)
	default:
		return "Не понимаю. Используй /help для списка команд."
	}
}

// === Command handlers ===

func (h *DevStudioHandler) cmdStart() string {
	return `Привет! ⚡ Добро пожаловать в Dev Studio!

Здесь ты можешь:
• Создать новое приложение
• Настроить команды и автоответы
• Получить API токен для интеграций
• Настроить вебхуки

Начни с /newapp чтобы создать своё первое приложение!

/help - все команды`
}

func (h *DevStudioHandler) cmdHelp() string {
	return `📋 Доступные команды:

🆕 Создание и управление
/newapp - Создать новое приложение
/myapps - Список моих приложений
/edit - Редактировать приложение
/delete - Удалить приложение

⚙️ Настройки приложения
/token - Получить/обновить API токен
/commands - Настроить команды
/webhook - Настроить вебхук

❌ /cancel - Отменить текущее действие`
}

func (h *DevStudioHandler) cmdNewApp(userID uint) string {
	h.setState(userID, StateAwaitingAppName, "")
	return `Отлично! Давай создадим новое приложение.

Как будет называться твоё приложение?
(Например: "Crypto Trading" или "NFT Gallery")`
}

func (h *DevStudioHandler) cmdMyApps(userID uint) string {
	var apps []models.MiniApp
	h.db.Where("creator_id = ?", userID).Order("created_at DESC").Find(&apps)

	if len(apps) == 0 {
		return "У тебя пока нет приложений. Создай первое с помощью /newapp"
	}

	var sb strings.Builder
	sb.WriteString("📱 Твои приложения:\n\n")

	for i, app := range apps {
		status := "⏳ На модерации"
		if app.ModerationStatus == models.ModerationApproved {
			status = "✅ Активно"
		} else if app.ModerationStatus == models.ModerationRejected {
			status = "❌ Отклонено"
		}

		username := "не задан"
		if app.BotUsername != "" {
			username = "@" + app.BotUsername
		}

		sb.WriteString(fmt.Sprintf("%d. %s %s\n", i+1, app.Icon, app.Title))
		sb.WriteString(fmt.Sprintf("   Username: %s\n", username))
		sb.WriteString(fmt.Sprintf("   Статус: %s\n", status))
		sb.WriteString(fmt.Sprintf("   Пользователей: %d\n\n", app.UsersCount))
	}

	sb.WriteString("Используй /edit для редактирования или /token для получения токена")
	return sb.String()
}

func (h *DevStudioHandler) cmdToken(userID uint) string {
	var apps []models.MiniApp
	h.db.Where("creator_id = ?", userID).Find(&apps)

	if len(apps) == 0 {
		return "У тебя нет приложений. Создай первое с /newapp"
	}

	if len(apps) == 1 {
		return h.showToken(apps[0])
	}

	h.setState(userID, StateSelectingApp, `{"action":"token"}`)
	return h.formatAppList(apps, "Выбери приложение для получения токена (введи номер):")
}

func (h *DevStudioHandler) cmdEditApp(userID uint) string {
	var apps []models.MiniApp
	h.db.Where("creator_id = ?", userID).Find(&apps)

	if len(apps) == 0 {
		return "У тебя нет приложений. Создай первое с /newapp"
	}

	if len(apps) == 1 {
		h.setState(userID, StateEditingApp, fmt.Sprintf(`{"app_id":%d}`, apps[0].ID))
		return h.showEditMenu(apps[0])
	}

	h.setState(userID, StateSelectingApp, `{"action":"edit"}`)
	return h.formatAppList(apps, "Выбери приложение для редактирования (введи номер):")
}

func (h *DevStudioHandler) cmdDeleteApp(userID uint) string {
	var apps []models.MiniApp
	h.db.Where("creator_id = ?", userID).Find(&apps)

	if len(apps) == 0 {
		return "У тебя нет приложений."
	}

	if len(apps) == 1 {
		h.setState(userID, StateDeletingApp, fmt.Sprintf(`{"app_id":%d}`, apps[0].ID))
		return fmt.Sprintf("⚠️ Ты уверен, что хочешь удалить приложение \"%s\"?\n\nВведи YES для подтверждения или /cancel для отмены.", apps[0].Title)
	}

	h.setState(userID, StateSelectingApp, `{"action":"delete"}`)
	return h.formatAppList(apps, "Выбери приложение для удаления (введи номер):")
}

func (h *DevStudioHandler) cmdCommands(userID uint) string {
	var apps []models.MiniApp
	h.db.Where("creator_id = ?", userID).Find(&apps)

	if len(apps) == 0 {
		return "У тебя нет приложений. Создай первое с /newapp"
	}

	if len(apps) == 1 {
		h.setState(userID, StateAwaitingCommand, fmt.Sprintf(`{"app_id":%d}`, apps[0].ID))
		return h.showCommandsMenu(apps[0])
	}

	h.setState(userID, StateSelectingApp, `{"action":"commands"}`)
	return h.formatAppList(apps, "Выбери приложение для настройки команд (введи номер):")
}

func (h *DevStudioHandler) cmdWebhook(userID uint) string {
	var apps []models.MiniApp
	h.db.Where("creator_id = ?", userID).Find(&apps)

	if len(apps) == 0 {
		return "У тебя нет приложений. Создай первое с /newapp"
	}

	if len(apps) == 1 {
		h.setState(userID, StateAwaitingWebhook, fmt.Sprintf(`{"app_id":%d}`, apps[0].ID))
		currentWebhook := "не установлен"
		if apps[0].WebhookURL != "" {
			currentWebhook = apps[0].WebhookURL
		}
		return fmt.Sprintf("🔗 Вебхук для %s\n\nТекущий URL: %s\n\nВведи новый URL вебхука или /cancel для отмены:", apps[0].Title, currentWebhook)
	}

	h.setState(userID, StateSelectingApp, `{"action":"webhook"}`)
	return h.formatAppList(apps, "Выбери приложение для настройки вебхука (введи номер):")
}

// === State handlers ===

func (h *DevStudioHandler) handleAppName(userID uint, name string) string {
	if len(name) < 3 {
		return "Название слишком короткое. Минимум 3 символа."
	}
	if len(name) > 50 {
		return "Название слишком длинное. Максимум 50 символов."
	}

	h.setState(userID, StateAwaitingAppDesc, fmt.Sprintf(`{"name":"%s"}`, name))
	return fmt.Sprintf("Отличное название: %s!\n\nТеперь опиши, что делает твоё приложение (краткое описание):", name)
}

func (h *DevStudioHandler) handleAppDesc(userID uint, desc string) string {
	if len(desc) < 10 {
		return "Описание слишком короткое. Минимум 10 символов."
	}

	data := h.getStateData(userID)
	data["description"] = desc
	h.setStateWithData(userID, StateAwaitingAppIcon, data)

	return "Теперь выбери иконку (эмодзи) для приложения:\n\n(Отправь один эмодзи, например: 🤖 🎮 💰 🔥 ⚡)"
}

func (h *DevStudioHandler) handleAppIcon(userID uint, icon string) string {
	// Simple emoji validation - just check it's short
	if len(icon) > 10 {
		return "Отправь только один эмодзи"
	}

	data := h.getStateData(userID)
	data["icon"] = icon
	h.setStateWithData(userID, StateAwaitingCategory, data)

	// Get categories
	var categories []models.Category
	h.db.Order("\"order\" ASC").Find(&categories)

	var sb strings.Builder
	sb.WriteString("Выбери категорию для приложения (введи номер):\n\n")
	for i, cat := range categories {
		sb.WriteString(fmt.Sprintf("%d. %s %s\n", i+1, cat.Icon, cat.Name))
	}

	return sb.String()
}

func (h *DevStudioHandler) handleCategory(userID uint, input string) string {
	num, err := strconv.Atoi(input)
	if err != nil {
		return "Введи номер категории"
	}

	var categories []models.Category
	h.db.Order("\"order\" ASC").Find(&categories)

	if num < 1 || num > len(categories) {
		return "Неверный номер категории"
	}

	category := categories[num-1]
	data := h.getStateData(userID)
	data["category_id"] = float64(category.ID)
	h.setStateWithData(userID, StateAwaitingUsername, data)

	return fmt.Sprintf("Категория: %s %s\n\nТеперь придумай username для приложения.\n\nUsername должен:\n• Быть уникальным\n• Содержать только a-z, 0-9 и _\n• Заканчиваться на 'app'\n\n(Например: mytradingapp, nft_gallery_app)", category.Icon, category.Name)
}

func (h *DevStudioHandler) handleUsername(userID uint, username string) string {
	username = strings.ToLower(strings.TrimPrefix(username, "@"))

	if len(username) < 5 {
		return "Username слишком короткий. Минимум 5 символов."
	}
	if !strings.HasSuffix(username, "app") {
		return "Username должен заканчиваться на 'app'"
	}

	// Check uniqueness
	var existingApp models.MiniApp
	if err := h.db.Where("bot_username = ?", username).First(&existingApp).Error; err == nil {
		return "Этот username уже занят. Попробуй другой."
	}

	data := h.getStateData(userID)
	data["username"] = username
	h.setStateWithData(userID, StateAwaitingWelcome, data)

	return fmt.Sprintf("Username @%s свободен! ✅\n\nТеперь напиши приветственное сообщение.\n\nЭто сообщение увидят пользователи когда отправят /start:\n\n(Или отправь /skip чтобы пропустить)", username)
}

func (h *DevStudioHandler) handleWelcome(userID uint, welcome string) string {
	data := h.getStateData(userID)

	if welcome != "/skip" {
		data["welcome"] = welcome
	}

	// Create the app
	name, _ := data["name"].(string)
	desc, _ := data["description"].(string)
	icon, _ := data["icon"].(string)
	categoryID := uint(data["category_id"].(float64))
	username, _ := data["username"].(string)
	welcomeMsg, _ := data["welcome"].(string)

	apiToken := models.GenerateAPIToken()

	app := models.MiniApp{
		Title:            name,
		Subtitle:         desc,
		Description:      desc,
		Icon:             icon,
		CategoryID:       categoryID,
		CreatorID:        userID,
		BotUsername:      username,
		WelcomeMessage:   welcomeMsg,
		APIToken:         apiToken,
		ModerationStatus: models.ModerationPending,
		IsVerified:       false,
		IsSecret:         false,
		UsersCount:       0,
	}

	if err := h.db.Create(&app).Error; err != nil {
		h.setState(userID, StateIdle, "")
		return "Ошибка при создании приложения. Попробуй снова с /newapp"
	}

	// Create /start command if welcome message provided
	if welcomeMsg != "" {
		startCmd := models.BotCommand{
			AppID:       app.ID,
			Command:     "/start",
			Description: "Запустить приложение",
			Response:    welcomeMsg,
			IsEnabled:   true,
		}
		h.db.Create(&startCmd)
	}

	h.setState(userID, StateIdle, "")

	return fmt.Sprintf(`🎉 Поздравляю! Приложение создано!

%s %s
@%s

🔑 API Token (сохрани его!):
%s

📋 Статус: ⏳ На модерации
Приложение будет проверено в течение 24 часов.

Что дальше:
• /commands - добавить команды
• /webhook - настроить вебхук
• /token - показать токен снова

Документация API: /help`, app.Icon, app.Title, app.BotUsername, apiToken)
}

func (h *DevStudioHandler) handleAppSelection(userID uint, input string) string {
	num, err := strconv.Atoi(input)
	if err != nil {
		return "Введи номер приложения"
	}

	var apps []models.MiniApp
	h.db.Where("creator_id = ?", userID).Find(&apps)

	if num < 1 || num > len(apps) {
		return "Неверный номер"
	}

	app := apps[num-1]
	data := h.getStateData(userID)
	action, _ := data["action"].(string)

	switch action {
	case "token":
		h.setState(userID, StateIdle, "")
		return h.showToken(app)
	case "edit":
		h.setState(userID, StateEditingApp, fmt.Sprintf(`{"app_id":%d}`, app.ID))
		return h.showEditMenu(app)
	case "delete":
		h.setState(userID, StateDeletingApp, fmt.Sprintf(`{"app_id":%d}`, app.ID))
		return fmt.Sprintf("⚠️ Ты уверен, что хочешь удалить приложение \"%s\"?\n\nВведи YES для подтверждения или /cancel для отмены.", app.Title)
	case "commands":
		h.setState(userID, StateAwaitingCommand, fmt.Sprintf(`{"app_id":%d}`, app.ID))
		return h.showCommandsMenu(app)
	case "webhook":
		h.setState(userID, StateAwaitingWebhook, fmt.Sprintf(`{"app_id":%d}`, app.ID))
		currentWebhook := "не установлен"
		if app.WebhookURL != "" {
			currentWebhook = app.WebhookURL
		}
		return fmt.Sprintf("🔗 Вебхук для %s\n\nТекущий URL: %s\n\nВведи новый URL вебхука:", app.Title, currentWebhook)
	}

	return "Ошибка"
}

func (h *DevStudioHandler) handleEditChoice(userID uint, choice string) string {
	data := h.getStateData(userID)
	appID := uint(data["app_id"].(float64))

	var app models.MiniApp
	h.db.First(&app, appID)

	switch choice {
	case "1": // Edit name
		h.setState(userID, StateAwaitingNewName, fmt.Sprintf(`{"app_id":%d}`, appID))
		return fmt.Sprintf("Текущее название: %s\n\nВведи новое название:", app.Title)
	case "2": // Edit description
		h.setState(userID, StateAwaitingNewDesc, fmt.Sprintf(`{"app_id":%d}`, appID))
		return fmt.Sprintf("Текущее описание: %s\n\nВведи новое описание:", app.Description)
	case "3": // Commands
		h.setState(userID, StateAwaitingCommand, fmt.Sprintf(`{"app_id":%d}`, appID))
		return h.showCommandsMenu(app)
	case "4": // Webhook
		h.setState(userID, StateAwaitingWebhook, fmt.Sprintf(`{"app_id":%d}`, appID))
		currentWebhook := "не установлен"
		if app.WebhookURL != "" {
			currentWebhook = app.WebhookURL
		}
		return fmt.Sprintf("Текущий вебхук: %s\n\nВведи новый URL:", currentWebhook)
	case "5": // Token
		h.setState(userID, StateIdle, "")
		return h.showToken(app)
	default:
		return "Выбери пункт от 1 до 5"
	}
}

func (h *DevStudioHandler) handleNewName(userID uint, name string) string {
	if len(name) < 3 || len(name) > 50 {
		return "Название должно быть от 3 до 50 символов"
	}

	data := h.getStateData(userID)
	appID := uint(data["app_id"].(float64))

	h.db.Model(&models.MiniApp{}).Where("id = ?", appID).Updates(map[string]interface{}{
		"title":             name,
		"moderation_status": models.ModerationPending,
	})

	h.setState(userID, StateIdle, "")
	return fmt.Sprintf("✅ Название изменено на: %s\n\nПриложение отправлено на повторную модерацию.", name)
}

func (h *DevStudioHandler) handleNewDesc(userID uint, desc string) string {
	if len(desc) < 10 {
		return "Описание должно быть минимум 10 символов"
	}

	data := h.getStateData(userID)
	appID := uint(data["app_id"].(float64))

	h.db.Model(&models.MiniApp{}).Where("id = ?", appID).Updates(map[string]interface{}{
		"description":       desc,
		"subtitle":          desc,
		"moderation_status": models.ModerationPending,
	})

	h.setState(userID, StateIdle, "")
	return "✅ Описание обновлено!\n\nПриложение отправлено на повторную модерацию."
}

func (h *DevStudioHandler) handleNewCommand(userID uint, input string) string {
	data := h.getStateData(userID)
	appID := uint(data["app_id"].(float64))

	// Check for "delete X" command
	if strings.HasPrefix(input, "delete ") {
		cmdToDelete := strings.TrimPrefix(input, "delete ")
		if !strings.HasPrefix(cmdToDelete, "/") {
			cmdToDelete = "/" + cmdToDelete
		}
		h.db.Where("app_id = ? AND command = ?", appID, cmdToDelete).Delete(&models.BotCommand{})

		var app models.MiniApp
		h.db.First(&app, appID)
		return fmt.Sprintf("✅ Команда %s удалена!\n\n%s", cmdToDelete, h.showCommandsMenu(app))
	}

	if !strings.HasPrefix(input, "/") {
		return "Команда должна начинаться с /"
	}

	// Check if command exists
	var existingCmd models.BotCommand
	if err := h.db.Where("app_id = ? AND command = ?", appID, input).First(&existingCmd).Error; err == nil {
		return fmt.Sprintf("Команда %s уже существует. Для удаления напиши: delete %s", input, input)
	}

	data["new_command"] = input
	h.setStateWithData(userID, StateAwaitingCmdDesc, data)

	return fmt.Sprintf("Команда: %s\n\nВведи краткое описание команды:", input)
}

func (h *DevStudioHandler) handleCmdDescription(userID uint, desc string) string {
	data := h.getStateData(userID)
	data["cmd_desc"] = desc
	h.setStateWithData(userID, StateAwaitingCmdResp, data)

	return "Теперь введи ответ приложения на эту команду:"
}

func (h *DevStudioHandler) handleCmdResponse(userID uint, response string) string {
	data := h.getStateData(userID)
	appID := uint(data["app_id"].(float64))
	command, _ := data["new_command"].(string)
	description, _ := data["cmd_desc"].(string)

	cmd := models.BotCommand{
		AppID:       appID,
		Command:     command,
		Description: description,
		Response:    response,
		IsEnabled:   true,
	}
	h.db.Create(&cmd)

	var app models.MiniApp
	h.db.First(&app, appID)

	h.setState(userID, StateAwaitingCommand, fmt.Sprintf(`{"app_id":%d}`, appID))
	return fmt.Sprintf("✅ Команда %s добавлена!\n\n%s", command, h.showCommandsMenu(app))
}

func (h *DevStudioHandler) handleWebhookUrl(userID uint, url string) string {
	data := h.getStateData(userID)
	appID := uint(data["app_id"].(float64))

	if url == "/clear" || url == "clear" {
		h.db.Model(&models.MiniApp{}).Where("id = ?", appID).Update("webhook_url", "")
		h.setState(userID, StateIdle, "")
		return "✅ Вебхук удалён"
	}

	if !strings.HasPrefix(url, "https://") {
		return "URL должен начинаться с https://"
	}

	h.db.Model(&models.MiniApp{}).Where("id = ?", appID).Update("webhook_url", url)
	h.setState(userID, StateIdle, "")

	return fmt.Sprintf("✅ Вебхук установлен:\n%s\n\nТеперь сообщения пользователей будут отправляться на этот URL.", url)
}

func (h *DevStudioHandler) handleDeleteConfirm(userID uint, confirm string) string {
	if strings.ToUpper(confirm) != "YES" {
		h.setState(userID, StateIdle, "")
		return "Удаление отменено."
	}

	data := h.getStateData(userID)
	appID := uint(data["app_id"].(float64))

	var app models.MiniApp
	h.db.First(&app, appID)
	title := app.Title

	// Delete related data
	h.db.Where("app_id = ?", appID).Delete(&models.BotCommand{})
	h.db.Where("app_id = ?", appID).Delete(&models.AppMessage{})
	h.db.Where("app_id = ?", appID).Delete(&models.AppUser{})
	h.db.Where("app_id = ?", appID).Delete(&models.WebhookLog{})
	h.db.Delete(&app)

	h.setState(userID, StateIdle, "")
	return fmt.Sprintf("✅ Приложение \"%s\" удалено.", title)
}

// === Helper methods ===

func (h *DevStudioHandler) getState(userID uint) *models.ConversationState {
	var state models.ConversationState
	if err := h.db.Where("user_id = ?", userID).First(&state).Error; err != nil {
		state = models.ConversationState{
			UserID: userID,
			State:  StateIdle,
			Data:   "{}",
		}
		h.db.Create(&state)
	}
	return &state
}

func (h *DevStudioHandler) setState(userID uint, state string, data string) {
	if data == "" {
		data = "{}"
	}
	h.db.Model(&models.ConversationState{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"state":      state,
		"data":       data,
		"updated_at": time.Now(),
	})
}

func (h *DevStudioHandler) getStateData(userID uint) map[string]interface{} {
	state := h.getState(userID)
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(state.Data), &data); err != nil {
		return make(map[string]interface{})
	}
	return data
}

func (h *DevStudioHandler) setStateWithData(userID uint, stateName string, data map[string]interface{}) {
	jsonData, _ := json.Marshal(data)
	h.setState(userID, stateName, string(jsonData))
}

func (h *DevStudioHandler) formatAppList(apps []models.MiniApp, header string) string {
	var sb strings.Builder
	sb.WriteString(header + "\n\n")
	for i, app := range apps {
		sb.WriteString(fmt.Sprintf("%d. %s %s (@%s)\n", i+1, app.Icon, app.Title, app.BotUsername))
	}
	return sb.String()
}

func (h *DevStudioHandler) showToken(app models.MiniApp) string {
	return fmt.Sprintf(`🔑 API Token для %s %s

%s

⚠️ Храни токен в безопасности!

Для перегенерации токена используй /edit → Токен`, app.Icon, app.Title, app.APIToken)
}

func (h *DevStudioHandler) showEditMenu(app models.MiniApp) string {
	return fmt.Sprintf(`⚙️ Редактирование: %s %s

Выбери что изменить:

1. 📝 Название
2. 📄 Описание
3. 📋 Команды
4. 🔗 Вебхук
5. 🔑 Показать токен

Введи номер или /cancel для отмены:`, app.Icon, app.Title)
}

func (h *DevStudioHandler) showCommandsMenu(app models.MiniApp) string {
	var commands []models.BotCommand
	h.db.Where("app_id = ?", app.ID).Order("command ASC").Find(&commands)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 Команды приложения %s:\n\n", app.Title))

	if len(commands) == 0 {
		sb.WriteString("Пока нет команд.\n\n")
	} else {
		for _, cmd := range commands {
			sb.WriteString(fmt.Sprintf("%s - %s\n", cmd.Command, cmd.Description))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("Введи новую команду (например /help)\nИли 'delete /команда' для удаления\n\n/cancel - выход")

	return sb.String()
}
