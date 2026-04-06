package mailer

type Mailer interface {
	SendWelcomeEmail(to string) error
}
