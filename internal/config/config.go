package config

import (
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
}

func Load() *Config {
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	otpExpiry, _ := strconv.Atoi(getEnv("OTP_EXPIRY_MINUTES", "10"))
	rateLimitReqs, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))
	rateLimitWindow, _ := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW", "60"))

	return &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", "change-this-secret-key"),
		SMTPHost:           getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:           smtpPort,
		SMTPUser:           getEnv("SMTP_USER", ""),
		SMTPPassword:       getEnv("SMTP_PASSWORD", ""),
		OTPExpiryMinutes:   otpExpiry,
		RateLimitRequests:  rateLimitReqs,
		RateLimitWindow:    rateLimitWindow,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
