package email

import (
	"bytes"
	"fmt"
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

func NewEmailSender(config *SMTPConfig) *EmailSender {
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
	message := s.buildMessage(to, subject, body.String())

	// SMTPサーバーに接続して送信
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	
	// MailCatcherは認証不要
	err = smtp.SendMail(addr, nil, s.config.From, to, []byte(message))
	if err != nil {
		return fmt.Errorf("メール送信エラー: %w", err)
	}

	return nil
}

// メールメッセージを構築
func (s *EmailSender) buildMessage(to []string, subject string, htmlBody string) string {
	headers := make(map[string]string)
	headers["From"] = s.config.From
	headers["To"] = to[0] // 簡略化のため最初の宛先のみ
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	return message
}

func (s *EmailSender) SendSimpleEmail(to []string, subject string, body string) error {
	message := s.buildMessage(to, subject, body)
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	
	err := smtp.SendMail(addr, nil, s.config.From, to, []byte(message))
	if err != nil {
		return fmt.Errorf("メール送信エラー: %w", err)
	}

	return nil
}

func GetTemplateDir() string {
	return filepath.Join(".", "infrastructure", "emails", "templates")
}
