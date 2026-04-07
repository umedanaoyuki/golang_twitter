package mailcatcher

import (
	"golang_twitter/mailer"
	"log"
)

type MailCatcherMailer struct {
	// sender *email.EmailSender
}

// MailCatcher 用の Mailer を作成
func NewMailCatcherMailer() mailer.Mailer {
	return &MailCatcherMailer{}
}

// ウェルカムメールを送信
func (m *MailCatcherMailer) SendWelcomeEmail(to string) error {
	// templatePath := filepath.Join(email.GetTemplateDir(), "welcome.html")

	// data := map[string]interface{}{
	// 	"Email":   to,
	// 	"AppName": "Twitter Clone",
	// }

	// return m.sender.SendEmail(
	// 	[]string{to},
	// 	"ようこそ Twitter Clone へ",
	// 	templatePath,
	// 	data,
	// )
	log.Println("Welcome email 送られました")
	return nil
}
