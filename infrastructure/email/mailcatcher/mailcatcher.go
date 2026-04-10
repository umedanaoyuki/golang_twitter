package mailcatcher

import (
	"bytes"
	"fmt"
	"golang_twitter/mailer"
	"html/template"
	"log"
	"net/smtp"
	"path/filepath"
)

// MailCatcherConfig は MailCatcher の設定
type MailCatcherConfig struct {
	Host string
	Port int
	From string
}

// MailCatcherMailer は MailCatcher を使用したメール送信の実装
type MailCatcherMailer struct {
	config *MailCatcherConfig
}

// NewMailCatcherMailer は MailCatcher 用の Mailer を作成
func NewMailCatcherMailer() mailer.Mailer {
	config := &MailCatcherConfig{
		Host: "mail", // docker-composeのサービス名
		Port: 1025,
		From: "test@email.com",
	}
	return &MailCatcherMailer{config: config}
}

// SendWelcomeEmail はウェルカムメールを送信
func (m *MailCatcherMailer) SendWelcomeEmail(to string, token string) error {
	// テンプレートパス
	templatePath := filepath.Join(mailer.GetTemplateDir(), "welcome_stg.html")

	// テンプレートを読み込み
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("テンプレート読み込みエラー: %w", err)
	}

	// アクティベーションURLを生成
	activationURL := fmt.Sprintf("http://localhost:8080/activate?token=%s", token)

	// データ準備
	data := map[string]interface{}{
		"Email":         to,
		"ActivationURL": activationURL,
	}

	// HTMLを生成
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("テンプレート実行エラー: %w", err)
	}

	// メールメッセージを構築
	message := mailer.BuildMessage([]string{to}, "ようこそ Twitter Clone へ", body.String(), m.config.From)

	// SMTPサーバーに接続して送信
	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)

	// MailCatcherは認証不要
	err = smtp.SendMail(addr, nil, m.config.From, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("メール送信エラー: %w", err)
	}

	log.Printf("Welcome email sent to: %s", to)
	return nil
}
