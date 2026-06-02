package smtp

import (
	"bytes"
	"fmt"
	"golang_twitter/mailer"
	"html/template"
	"net/smtp"
	"path/filepath"
)

// SMTPConfig はメールサーバーの設定
type SMTPConfig struct {
	Host     string
	Port     int
	From     string
	Username string
	Password string
}

func NewMailCatcherConfig() *SMTPConfig {
	return &SMTPConfig{
		Host:     "mail", // docker-composeのサービス名
		Port:     1025,
		From:     "test@email.com",
		Username: "",
		Password: "",
	}
}
type EmailSender struct {
	config *SMTPConfig
}

func (m *EmailSender) SendWelcomeEmail(to string, token string) error {
	templatePath := filepath.Join(mailer.GetTemplateDir(), "welcome.html")

	data := map[string]interface{}{
		"Email":   to,
		"Token":   token,
		"AppName": "Twitter Clone",
	}

	return m.SendEmail(
		[]string{to},
		"ようこそ Twitter Clone へ",
		templatePath,
		data,
	)
}

func NewEmailSender(config *SMTPConfig) mailer.Mailer {
	return &EmailSender{config: config}
}

// メールを送信
func (s *EmailSender) SendEmail(to []string, subject string, templatePath string, data interface{}) error {
	// テンプレートを読み込み
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("テンプレート読み込みエラー: %w", err)
	}

	// HTMLを生成
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("テンプレート実行エラー: %w", err)
	}

	// メールメッセージを構築
	message := mailer.BuildMessage(to, subject, body.String(), s.config.From)

	// SMTPサーバーに接続して送信
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	
	// MailCatcherは認証不要
	err = smtp.SendMail(addr, nil, s.config.From, to, []byte(message))
	if err != nil {
		return fmt.Errorf("メール送信エラー: %w", err)
	}

	return nil
}
