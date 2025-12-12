package database

import (
	"github.com/fasad/solanafon-back/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// Users & Auth
		&models.User{},
		&models.OTP{},
		&models.ManaTransaction{},

		// Apps
		&models.Category{},
		&models.MiniApp{},
		&models.AppUser{},
		&models.AppMessage{},

		// Bot system
		&models.BotCommand{},
		&models.WebhookLog{},
		&models.ConversationState{},

		// Secret Login
		&models.SecretNumber{},
		&models.SecretAccess{},
	)
}
