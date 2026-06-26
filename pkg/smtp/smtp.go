package smtp

import (
	"log/slog"
	"net/smtp"
	"sso/internal/config"
)

func SendEmailNotification(subject string, body string, to []string, cfg *config.Config, log *slog.Logger) error {
	auth := smtp.PlainAuth(
		"",
		cfg.SMTP.Email,
		cfg.SMTP.Password,
		"smtp.gmail.com",
	)
	msg := "Subject: " + subject + "\n" + body
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		cfg.SMTP.Email,
		to,
		[]byte(msg),
	)
	if err != nil {
		log.Error("error while sending email", "error", err.Error())
		return err
	}
	return nil
}
