package service

import (
	"banksystem/internal/config"
	"fmt"
	"strconv"

	"gopkg.in/mail.v2"
)

type SMTPService struct {
	config *config.Config
}

func NewSMTPService(host, port, username, password string) *SMTPService {
	return &SMTPService{
		config: &config.Config{
			SMTPHost:     host,
			SMTPPort:     port,
			SMTPUsername: username,
			SMTPPassword: password,
		},
	}
}

func (s *SMTPService) SendEmail(to, subject, body string) error {
	port, err := strconv.Atoi(s.config.SMTPPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	m := mail.NewMessage()
	m.SetHeader("From", s.config.SMTPUsername)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := mail.NewDialer(s.config.SMTPHost, port, s.config.SMTPUsername, s.config.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (s *SMTPService) SendLowBalanceNotification(email string, balance float64) error {
	subject := "Low Balance Alert"
	body := fmt.Sprintf(`
		<h1>Low Balance Alert</h1>
		<p>Your account balance is now %.2f</p>
		<p>Please consider adding funds to your account.</p>
	`, balance)

	return s.SendEmail(email, subject, body)
}

func (s *SMTPService) SendTransactionNotification(email string, amount float64, transactionType string) error {
	subject := "Transaction Notification"
	body := fmt.Sprintf(`
		<h1>Transaction Notification</h1>
		<p>A %s transaction of %.2f has been processed on your account.</p>
	`, transactionType, amount)

	return s.SendEmail(email, subject, body)
}

func (s *SMTPService) SendCreditApprovalNotification(email string, amount float64, term int) error {
	subject := "Credit Application Approved"
	body := fmt.Sprintf(`
		<h1>Credit Application Approved</h1>
		<p>Your credit application for %.2f with a term of %d months has been approved.</p>
	`, amount, term)

	return s.SendEmail(email, subject, body)
}

func (s *SMTPService) SendPaymentReminderNotification(email string, amount float64, dueDate string) error {
	subject := "Payment Reminder"
	body := fmt.Sprintf(`
		<h1>Payment Reminder</h1>
		<p>You have a payment of %.2f due on %s.</p>
		<p>Please ensure you have sufficient funds in your account.</p>
	`, amount, dueDate)

	return s.SendEmail(email, subject, body)
}
