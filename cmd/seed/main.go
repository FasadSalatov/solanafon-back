package main

import (
	"log"

	"github.com/fasad/solanafon-back/internal/config"
	"github.com/fasad/solanafon-back/internal/database"
	"github.com/fasad/solanafon-back/internal/models"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load config
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Seeding database...")

	// ==================== CATEGORIES ====================
	// Based on screenshots: AI, Игры (Games), Трейдинг (Trading), DePIN, DeFi, NFT, Стейкинг (Staking), Сервисы (Services)
	categories := []models.Category{
		{Name: "AI", Slug: "ai", Description: "AI-powered applications", Icon: "🤖", Order: 1},
		{Name: "Games", Slug: "games", Description: "Play-to-earn games and entertainment", Icon: "🎮", Order: 2},
		{Name: "Trading", Slug: "trading", Description: "Trading and market analysis tools", Icon: "📊", Order: 3},
		{Name: "DePIN", Slug: "depin", Description: "Decentralized Physical Infrastructure", Icon: "🌐", Order: 4},
		{Name: "DeFi", Slug: "defi", Description: "Decentralized Finance applications", Icon: "💎", Order: 5},
		{Name: "NFT", Slug: "nft", Description: "NFT marketplaces and collections", Icon: "🖼️", Order: 6},
		{Name: "Staking", Slug: "staking", Description: "Staking and yield farming", Icon: "🔒", Order: 7},
		{Name: "Services", Slug: "services", Description: "Various utility services", Icon: "🛠️", Order: 8},
	}

	for _, category := range categories {
		db.FirstOrCreate(&category, models.Category{Slug: category.Slug})
	}
	log.Println("✓ Categories seeded")

	// Get category IDs
	var catAI, catGames, catTrading, catDePIN, catDeFi, catNFT, catStaking, catServices models.Category
	db.Where("slug = ?", "ai").First(&catAI)
	db.Where("slug = ?", "games").First(&catGames)
	db.Where("slug = ?", "trading").First(&catTrading)
	db.Where("slug = ?", "depin").First(&catDePIN)
	db.Where("slug = ?", "defi").First(&catDeFi)
	db.Where("slug = ?", "nft").First(&catNFT)
	db.Where("slug = ?", "staking").First(&catStaking)
	db.Where("slug = ?", "services").First(&catServices)

	// ==================== DEV STUDIO (System App) ====================
	// Dev Studio is the official app for creating and managing mini-apps
	devStudio := models.MiniApp{
		Title:       "Dev Studio",
		Subtitle:    "Создавай и управляй своими приложениями",
		Icon:        "⚡",
		CategoryID:  catServices.ID,
		IsSecret:    false,
		IsVerified:  true,
		UsersCount:  0,
		Description: "Dev Studio - официальное приложение для разработчиков. Создавай мини-приложения, настраивай команды, получай API токены для интеграций.",
		BotUsername: "devstudio",
		APIToken:    "SYSTEM_DEVSTUDIO_TOKEN", // Special system token
		WelcomeMessage: `Привет! ⚡ Добро пожаловать в Dev Studio!

Здесь ты можешь:
• Создать новое приложение
• Настроить команды и автоответы
• Получить API токен для интеграций
• Настроить вебхуки

Начни с /newapp чтобы создать своё первое приложение!

/help - все команды`,
		ModerationStatus: models.ModerationApproved,
	}
	db.FirstOrCreate(&devStudio, models.MiniApp{BotUsername: "devstudio"})
	log.Println("✓ Dev Studio seeded")

	// Dev Studio commands
	if devStudio.ID > 0 {
		devStudioCommands := []models.BotCommand{
			{AppID: devStudio.ID, Command: "/start", Description: "Начать работу", Response: devStudio.WelcomeMessage, IsEnabled: true},
			{AppID: devStudio.ID, Command: "/newapp", Description: "Создать новое приложение", Response: "", IsEnabled: true},
			{AppID: devStudio.ID, Command: "/myapps", Description: "Мои приложения", Response: "", IsEnabled: true},
			{AppID: devStudio.ID, Command: "/edit", Description: "Редактировать приложение", Response: "", IsEnabled: true},
			{AppID: devStudio.ID, Command: "/delete", Description: "Удалить приложение", Response: "", IsEnabled: true},
			{AppID: devStudio.ID, Command: "/token", Description: "Получить API токен", Response: "", IsEnabled: true},
			{AppID: devStudio.ID, Command: "/commands", Description: "Настроить команды", Response: "", IsEnabled: true},
			{AppID: devStudio.ID, Command: "/webhook", Description: "Настроить вебхук", Response: "", IsEnabled: true},
			{AppID: devStudio.ID, Command: "/help", Description: "Справка", Response: "", IsEnabled: true},
			{AppID: devStudio.ID, Command: "/cancel", Description: "Отменить действие", Response: "", IsEnabled: true},
		}

		for _, cmd := range devStudioCommands {
			db.FirstOrCreate(&cmd, models.BotCommand{AppID: cmd.AppID, Command: cmd.Command})
		}
		log.Println("✓ Dev Studio Commands seeded")
	}

	// ==================== SECRET NUMBERS ====================
	// Based on screenshots: +999 (XXX) XXX-XX-XX format, prices in MP
	secretNumbers := []models.SecretNumber{
		{Number: "+999 (764) 123-45-67", IsPremium: false, PriceMP: 50, IsAvailable: true},
		{Number: "+999 (831) 987-65-43", IsPremium: false, PriceMP: 50, IsAvailable: true},
		{Number: "+999 (555) 000-11-22", IsPremium: true, PriceMP: 100, IsAvailable: true},
		{Number: "+999 (123) 456-78-90", IsPremium: false, PriceMP: 50, IsAvailable: true},
		{Number: "+999 (999) 888-77-66", IsPremium: true, PriceMP: 100, IsAvailable: true},
		{Number: "+999 (777) 111-22-33", IsPremium: false, PriceMP: 50, IsAvailable: true},
		{Number: "+999 (333) 444-55-66", IsPremium: false, PriceMP: 50, IsAvailable: true},
		{Number: "+999 (100) 200-30-40", IsPremium: true, PriceMP: 150, IsAvailable: true},
	}

	for _, number := range secretNumbers {
		db.FirstOrCreate(&number, models.SecretNumber{Number: number.Number})
	}
	log.Println("✓ Secret Numbers seeded")

	log.Println("✅ Database seeding completed!")
	log.Println("")
	log.Println("Created:")
	log.Println("  - 8 categories")
	log.Println("  - 1 system app (Dev Studio)")
	log.Println("  - 8 secret numbers")
	log.Println("")
	log.Println("Users can now create their own apps via Dev Studio!")
}
