package email

import (
	"fin_notifications/internal/log"
	"net/smtp"
)

func NotifyByEmail(subject string, text string, to []string) error {
	// Настройки
	smtpHost := "smtp.mail.ru"         // SMTP сервер
	smtpPort := "587"                  // Порт SMTP
	email := "xp_89@mail.ru"           // Ваш email
	password := "u8dL9HXh6TAcxmhxJNWg" // Ваш пароль

	message := []byte("Subject: " + subject + "\r\n" + text)

	// Подключение к SMTP серверу
	auth := smtp.PlainAuth("", email, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, email, to, message)
	if err != nil {
		return err
	}

	log.Info("Email sent successfully!")

	return nil
}
