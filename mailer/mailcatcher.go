package mailer

import (
	"golang_twitter/infrastructure/email"
	"path/filepath"
)

type MailCatcherMailer struct {
	sender *email.EmailSender
}

// MailCatcher 用の Mailer を作成
func NewMailCatcherMailer() Mailer {
	config := email.NewMailCatcherConfig()
	sender := email.NewEmailSender(config)
	return &MailCatcherMailer{sender: sender}
}

// ウェルカムメールを送信
func (m *MailCatcherMailer) SendWelcomeEmail(to string) error {
	templatePath := filepath.Join(email.GetTemplateDir(), "welcome.html")

	data := map[string]interface{}{
		"Email":   to,
		"AppName": "Twitter Clone",
	}

	return m.sender.SendEmail(
		[]string{to},
		"ようこそ Twitter Clone へ",
		templatePath,
		data,
	)
}
