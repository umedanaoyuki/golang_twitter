package mailcatcher

import (
	"fmt"
	"log"
	"net/smtp"

	"golang_twitter/mailer"
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
	// アクティベーションURLを生成
	activationURL := fmt.Sprintf("http://localhost:8080/activate?token=%s", token)

	// HTMLを直接構築
	htmlBody := fmt.Sprintf(`
<!doctype html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>ようこそ</title>
  </head>
  <body>
    <h1>ようこそ Twitter Clone へ！</h1>
    <p>ご登録ありがとうございます。</p>
    <p>登録されたメールアドレス: <strong>%s</strong></p>
    <p>アカウントを有効化するには、以下のリンクをクリックしてください。</p>
    <a href="%s">アカウントを有効化</a>
    <hr />
    <p>このメールに心当たりがない場合は、無視してください。</p>
  </body>
</html>
	`, to, activationURL)

	// メールメッセージを構築
	message := mailer.BuildMessage([]string{to}, "ようこそ Twitter Clone へ", htmlBody, m.config.From)

	// SMTPサーバーに接続して送信
	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)

	// MailCatcherは認証不要
	err := smtp.SendMail(addr, nil, m.config.From, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("メール送信エラー: %w", err)
	}

	log.Printf("Welcome email sent to: %s", to)
	return nil
}
