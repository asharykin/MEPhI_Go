package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/go-mail/mail/v2"
)

type NotificationService struct {
	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPassword string
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		smtpHost:     os.Getenv("SMTP_HOST"),
		smtpPort:     587,
		smtpUser:     os.Getenv("SMTP_USER"),
		smtpPassword: os.Getenv("SMTP_PASSWORD"),
	}
}

func (s *NotificationService) SendPaymentNotification(email string, amount float64) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.smtpUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Уведомление о платеже")
	m.SetBody("text/html", fmt.Sprintf(`
		<h1>Уведомление о платеже</h1>
		<p>Сумма платежа: %.2f RUB</p>
		<p>Спасибо за использование нашего сервиса!</p>
	`, amount))

	d := mail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUser, s.smtpPassword)
	d.TLSConfig = &tls.Config{
		ServerName:         s.smtpHost,
		InsecureSkipVerify: false,
	}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Ошибка отправки email: %v", err)
		return err
	}

	return nil
}

func (s *NotificationService) SendCreditPaymentNotification(email string, amount float64, dueDate string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.smtpUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Напоминание о платеже по кредиту")
	m.SetBody("text/html", fmt.Sprintf(`
		<h1>Напоминание о платеже по кредиту</h1>
		<p>Сумма платежа: %.2f RUB</p>
		<p>Дата платежа: %s</p>
		<p>Пожалуйста, убедитесь, что на вашем счете достаточно средств.</p>
	`, amount, dueDate))

	d := mail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUser, s.smtpPassword)
	d.TLSConfig = &tls.Config{
		ServerName:         s.smtpHost,
		InsecureSkipVerify: false,
	}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Ошибка отправки email: %v", err)
		return err
	}

	return nil
}
