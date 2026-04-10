package mailer

import (
	"fmt"
	"path/filepath"
)

type Mailer interface {
	SendWelcomeEmail(to string, token string) error
}

func GetTemplateDir() string {
	return filepath.Join(".", "infrastructure", "emails", "templates")
}

// メールメッセージを構築
func BuildMessage(to []string, subject string, htmlBody string, from string) string {
	headers := make(map[string]string)
	headers["From"] = from
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