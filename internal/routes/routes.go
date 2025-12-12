package routes

import (
	"github.com/fasad/solanafon-back/internal/config"
	"github.com/fasad/solanafon-back/internal/handlers"
	"github.com/fasad/solanafon-back/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(api fiber.Router, db *gorm.DB, cfg *config.Config) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	miniAppHandler := handlers.NewMiniAppHandler(db)
	profileHandler := handlers.NewProfileHandler(db)
	secretHandler := handlers.NewSecretHandler(db)
	botHandler := handlers.NewBotHandler(db)
	devStudioHandler := handlers.NewDevStudioHandler(db)

	// Auth middleware
	authMiddleware := middleware.AuthRequired(cfg.JWTSecret)

	// ==================== PUBLIC ROUTES ====================

	// Authentication
	auth := api.Group("/auth")
	auth.Post("/email/request", authHandler.RequestOTP)
	auth.Post("/email/verify", authHandler.VerifyOTP)

	// ==================== PROTECTED ROUTES ====================

	// Auth - protected
	auth.Get("/me", authMiddleware, authHandler.GetMe)
	auth.Post("/logout", authMiddleware, authHandler.Logout)

	// Categories
	categories := api.Group("/categories", authMiddleware)
	categories.Get("/", miniAppHandler.GetCategories)
	categories.Get("/:slug/apps", miniAppHandler.GetByCategory)

	// Apps (MiniApps)
	apps := api.Group("/apps", authMiddleware)
	apps.Get("/", miniAppHandler.GetAll)             // List all apps (with optional ?category=slug filter)
	apps.Get("/search", miniAppHandler.Search)       // Search apps ?q=query
	apps.Post("/", miniAppHandler.CreateApp)         // Create new app (returns API token)
	apps.Get("/my", miniAppHandler.GetMyApps)        // Get user's own apps
	apps.Get("/:id", miniAppHandler.GetByID)         // Get app details
	apps.Put("/:id", miniAppHandler.UpdateApp)       // Update own app
	apps.Delete("/:id", miniAppHandler.DeleteApp)    // Delete own app

	// App Messages (chat)
	apps.Get("/:id/messages", miniAppHandler.GetAppMessages)   // Get chat history
	apps.Post("/:id/messages", miniAppHandler.SendAppMessage)  // Send message to app

	// Dev Studio special endpoint (for creating/managing apps via chat)
	apps.Post("/devstudio/message", devStudioHandler.ProcessMessage) // Dev Studio chat

	// App Settings (for app owner)
	apps.Get("/:id/settings", miniAppHandler.GetAppSettings)           // Get app settings with API token
	apps.Post("/:id/regenerate-token", miniAppHandler.RegenerateAPIToken) // Regenerate API token

	// Bot Commands (for app owner)
	apps.Get("/:id/commands", miniAppHandler.GetBotCommands)           // Get bot commands
	apps.Post("/:id/commands", miniAppHandler.AddBotCommand)           // Add bot command
	apps.Put("/:id/commands/:cmdId", miniAppHandler.UpdateBotCommand)  // Update bot command
	apps.Delete("/:id/commands/:cmdId", miniAppHandler.DeleteBotCommand) // Delete bot command

	// Profile
	profile := api.Group("/profile", authMiddleware)
	profile.Get("/", profileHandler.GetProfile)               // Get full profile with stats
	profile.Patch("/", profileHandler.UpdateProfile)          // Update profile info
	profile.Put("/settings", profileHandler.UpdateSettings)   // Update settings (notifications, 2FA)

	// Mana Points
	mana := api.Group("/mana", authMiddleware)
	mana.Get("/", profileHandler.GetManaHistory)    // Get MP balance and history
	mana.Post("/topup", profileHandler.TopUpMana)   // Top up MP (mock payment)

	// Secret Login
	secret := api.Group("/secret", authMiddleware)
	secret.Get("/numbers", secretHandler.GetNumbers)         // Get available numbers
	secret.Post("/activate", secretHandler.ActivateNumber)   // Activate (purchase) a number
	secret.Get("/status", secretHandler.GetStatus)           // Get secret access status
	secret.Delete("/deactivate", secretHandler.DeactivateSecret) // Deactivate secret access

	// ==================== BOT API (for external services) ====================
	// These endpoints use API token authentication (not JWT)
	bot := api.Group("/bot")
	bot.Get("/getMe", botHandler.GetMe)                    // Get bot info
	bot.Post("/sendMessage", botHandler.SendMessage)       // Send message to user
	bot.Get("/getUpdates", botHandler.GetUpdates)          // Get pending messages (polling)
	bot.Post("/setWebhook", botHandler.SetWebhook)         // Set webhook URL
	bot.Post("/deleteWebhook", botHandler.DeleteWebhook)   // Delete webhook
	bot.Get("/getWebhookInfo", botHandler.GetWebhookInfo)  // Get webhook info
	bot.Post("/setMyCommands", botHandler.SetCommands)     // Set bot commands
	bot.Get("/getMyCommands", botHandler.GetCommands)      // Get bot commands
}
