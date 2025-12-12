package utils

import (
	"fmt"
	"net/smtp"
)

type EmailConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

func SendOTPEmail(to, code string, config EmailConfig) error {
	from := config.User
	subject := "Your Solafon Login Code"
	body := fmt.Sprintf(`
Hello,

Your one-time password (OTP) for Solafon is: %s

This code will expire in 10 minutes.

If you didn't request this code, please ignore this email.

Best regards,
Solafon Team
`, code)

	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", from, to, subject, body))

	auth := smtp.PlainAuth("", config.User, config.Password, config.Host)
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	return smtp.SendMail(addr, auth, from, []string{to}, message)
}
