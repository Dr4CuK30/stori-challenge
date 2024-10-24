package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type GmailResource struct{}

func (r *GmailResource) SendEmail(destinatary string, subjectData string, body string) error {
	senderEmail := os.Getenv("SENDER_EMAIL")
	password := os.Getenv("APP_PASSWORD_SMTP")
	to := []string{destinatary}
	subject := subjectData
	message := []byte("From: " + senderEmail + "\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" + body + "\r\n")
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", senderEmail, password, "smtp.gmail.com"),
		senderEmail, []string{destinatary}, []byte(message))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}
	fmt.Println("Email sent successfully")
	return nil
}
