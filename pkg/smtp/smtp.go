package smtp

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"sso/internal/config"
)

func SendEmailNotification(subject string, body string, to []string, cfg config.SMTPConfig, log *slog.Logger) error {
	host := cfg.Host
	if host == "" {
		host = "smtp.gmail.com"
	}

	auth := smtp.PlainAuth(
		"",
		cfg.Email,
		cfg.Password,
		host,
	)

	msg := fmt.Sprintf("Subject: %s\nMIME-Version: 1.0\nContent-Type: text/plain; charset=UTF-8\n\n%s", subject, body)

	addr := fmt.Sprintf("%s:%d", host, cfg.Port)

	err := smtp.SendMail(
		addr,
		auth,
		cfg.Email,
		to,
		[]byte(msg),
	)
	if err != nil {
		log.Error("failed to send email", slog.String("error", err.Error()), slog.String("to", to[0]))
		return fmt.Errorf("smtp send error: %w", err)
	}
	return nil
}
