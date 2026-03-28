package routes

import (
	"github.com/fasad/solanafon-back/internal/config"
	"github.com/fasad/solanafon-back/internal/handlers"
	"github.com/fasad/solanafon-back/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupAPI registers all /api/* routes for the mobile app
func SetupAPI(api fiber.Router, db *gorm.DB, cfg *config.Config) {
	// Initialize handlers
	authV2 := handlers.NewAuthV2Handler(db, cfg)
	users := handlers.NewUsersHandler(db)
	convs := handlers.NewConversationsHandler(db)
	wallet := handlers.NewWalletHandler(db, cfg)
	notifications := handlers.NewNotificationsHandler(db)
	news := handlers.NewNewsHandler(db)
	calls := handlers.NewCallsHandler(db)
	referral := handlers.NewReferralHandler(db)
	support := handlers.NewSupportHandler(db)
	crash := handlers.NewCrashHandler(db)
	developer := handlers.NewDeveloperHandler(db, cfg)

	// Auth middleware
	auth := middleware.AuthRequired(cfg.JWTSecret)

	// ==================== AUTH (public) ====================
	authGroup := api.Group("/auth")
	authGroup.Post("/send-code", authV2.SendCode)
	authGroup.Post("/verify-code", authV2.VerifyCode)
	authGroup.Post("/refresh", authV2.RefreshToken)
	authGroup.Post("/logout", auth, authV2.LogoutV2)

	// ==================== USERS (protected) ====================
	usersGroup := api.Group("/users", auth)
	usersGroup.Get("/me", users.GetMe)
	usersGroup.Patch("/me", users.UpdateMe)
	usersGroup.Put("/me/settings", users.UpdateSettings)
	usersGroup.Delete("/me", users.DeleteMe)
	usersGroup.Get("/me/sessions", users.GetSessions)
	usersGroup.Delete("/me/sessions/:sessionId", users.RevokeSession)
	usersGroup.Delete("/me/sessions", users.RevokeAllSessions)

	// ==================== CONVERSATIONS (protected) ====================
	convsGroup := api.Group("/conversations", auth)
	convsGroup.Get("/", convs.ListConversations)
	convsGroup.Post("/", convs.StartConversation)
	convsGroup.Get("/:conversationId/messages", convs.GetMessages)
	convsGroup.Post("/:conversationId/messages", convs.SendMessage)
	convsGroup.Post("/:conversationId/messages/:messageId/callback", convs.ButtonCallback)
	convsGroup.Post("/:conversationId/read", convs.MarkAsRead)
	convsGroup.Delete("/:conversationId", convs.DeleteConversation)

	// ==================== APPS MARKETPLACE (protected) ====================
	appsGroup := api.Group("/apps", auth)
	appsGroup.Get("/", developer.ListApps)
	appsGroup.Get("/categories", developer.GetCategories)
	appsGroup.Get("/:appId", developer.GetAppDetail)
	appsGroup.Post("/:appId/launch", developer.LaunchApp)

	// ==================== DEVELOPER (protected) ====================
	devGroup := api.Group("/developer", auth)
	devGroup.Get("/apps", developer.ListMyApps)
	devGroup.Post("/apps", developer.CreateApp)
	devGroup.Get("/apps/:appId", developer.GetApp)
	devGroup.Put("/apps/:appId", developer.UpdateApp)
	devGroup.Delete("/apps/:appId", developer.DeleteApp)
	devGroup.Post("/apps/:appId/api-keys", developer.GenerateAPIKey)
	devGroup.Get("/apps/:appId/api-keys", developer.ListAPICredentials)
	devGroup.Delete("/apps/:appId/api-keys/:keyId", developer.RevokeAPIKey)
	devGroup.Put("/apps/:appId/webhook", developer.UpdateWebhook)
	devGroup.Get("/apps/:appId/welcome-message", developer.GetWelcomeMessage)
	devGroup.Put("/apps/:appId/welcome-message", developer.UpdateWelcomeMessage)

	// ==================== UPLOAD (protected) ====================
	api.Post("/upload", auth, developer.Upload)

	// ==================== WALLET (protected) ====================
	walletGroup := api.Group("/wallet", auth)
	walletGroup.Get("/balance", wallet.GetBalance)
	walletGroup.Post("/send", wallet.SendTransaction)
	walletGroup.Get("/transactions", wallet.GetTransactions)
	walletGroup.Get("/tokens", wallet.GetTokens)
	walletGroup.Get("/prices", wallet.GetPrices)
	walletGroup.Get("/transactions/:signature/status", wallet.GetTransactionStatus)

	// ==================== MANA (protected) ====================
	manaGroup := api.Group("/mana", auth)
	manaGroup.Get("/", wallet.GetManaPoints)
	manaGroup.Get("/tariffs", wallet.GetTariffs)
	manaGroup.Post("/purchase", wallet.PurchaseMana)
	manaGroup.Post("/gift", wallet.GiftMana)
	manaGroup.Get("/networks", wallet.GetNetworks)

	// ==================== NOTIFICATIONS (protected) ====================
	notifGroup := api.Group("/notifications", auth)
	notifGroup.Post("/push-token", notifications.RegisterPushToken)
	notifGroup.Delete("/push-token", notifications.UnregisterPushToken)
	notifGroup.Get("/", notifications.ListNotifications)
	notifGroup.Post("/:notificationId/read", notifications.MarkAsRead)
	notifGroup.Post("/read-all", notifications.MarkAllAsRead)
	notifGroup.Get("/unread-count", notifications.GetUnreadCount)

	// ==================== NEWS (protected) ====================
	newsGroup := api.Group("/news", auth)
	newsGroup.Get("/feed", news.GetFeed)
	newsGroup.Post("/:postId/like", news.LikePost)
	newsGroup.Post("/:postId/share", news.SharePost)
	newsGroup.Get("/:postId/comments", news.GetComments)
	newsGroup.Post("/:postId/comments", news.PostComment)
	newsGroup.Post("/", news.CreatePost)

	// ==================== CALLS (protected) ====================
	callsGroup := api.Group("/calls", auth)
	callsGroup.Post("/rooms/create", calls.CreateRoom)
	callsGroup.Get("/rooms/code/:roomCode", calls.GetRoomByCode)
	callsGroup.Post("/rooms/:roomId/join", calls.JoinRoom)
	callsGroup.Post("/rooms/:roomId/end", calls.EndCall)
	callsGroup.Patch("/rooms/:roomId/status", calls.UpdateParticipantStatus)

	// ==================== REFERRAL (protected) ====================
	referralGroup := api.Group("/referral", auth)
	referralGroup.Get("/", referral.GetReferralInfo)
	referralGroup.Get("/validate/:code", referral.ValidateCode)

	// ==================== SUPPORT (public/mixed) ====================
	supportGroup := api.Group("/support")
	supportGroup.Get("/faq", support.GetFAQ)
	supportGroup.Post("/tickets", auth, support.CreateTicket)

	// ==================== LEGAL (public) ====================
	legalGroup := api.Group("/legal")
	legalGroup.Get("/terms", support.GetTerms)
	legalGroup.Get("/privacy", support.GetPrivacy)

	// ==================== I18N (public) ====================
	api.Get("/i18n/languages", support.GetLanguages)

	// ==================== CRASH REPORTING (public) ====================
	api.Post("/crash/mobile", crash.ReportCrash)
}
