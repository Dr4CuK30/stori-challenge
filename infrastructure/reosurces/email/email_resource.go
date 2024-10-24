package email

type EmailResource interface {
	SendEmail(to string, subject string, body string) error
}
