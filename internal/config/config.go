package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	SMTPHost           string
	SMTPPort           int
	SMTPUser           string
	SMTPPassword       string
	OTPExpiryMinutes   int
	RateLimitRequests  int
	RateLimitWindow    int
	SolanaRPCURL       string
	UploadDir          string
	BaseURL            string
}

func Load() *Config {
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "465"))
	otpExpiry, _ := strconv.Atoi(getEnv("OTP_EXPIRY_MINUTES", "10"))
	rateLimitReqs, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))
	rateLimitWindow, _ := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW", "60"))

	// Support both DATABASE_URL and individual DB_* env vars
	dbURL := getEnv("DATABASE_URL", "")
	if dbURL == "" {
		dbHost := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "5432")
		dbUser := getEnv("DB_USER", "solafon")
		dbPass := getEnv("DB_PASSWORD", "")
		dbName := getEnv("DB_NAME", "solafon")
		dbSSL := getEnv("DB_SSLMODE", "disable")
		dbURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPass, dbName, dbSSL)
	}

	return &Config{
		Port:               getEnv("SERVER_PORT", getEnv("PORT", "8080")),
		DatabaseURL:        dbURL,
		JWTSecret:          getEnv("JWT_SECRET", "change-this-secret-key"),
		SMTPHost:           getEnv("SMTP_HOST", "smtp0001.neo.space"),
		SMTPPort:           smtpPort,
		SMTPUser:           getEnv("SMTP_USER", "Info@solafon.com"),
		SMTPPassword:       getEnv("SMTP_PASSWORD", ""),
		OTPExpiryMinutes:   otpExpiry,
		RateLimitRequests:  rateLimitReqs,
		RateLimitWindow:    rateLimitWindow,
		SolanaRPCURL:       getEnv("SOLANA_RPC_URL", "https://api.mainnet-beta.solana.com"),
		UploadDir:          getEnv("STORAGE_PATH", getEnv("UPLOAD_DIR", "./uploads")),
		BaseURL:            getEnv("APP_URL", getEnv("BASE_URL", "https://api.solafon.com")),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
