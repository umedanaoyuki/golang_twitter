package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	MaxImageSize = 5 << 20 // 5MB
)

var allowedImageMIMEs = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// AllowedImageMIME は対応する Content-Type かどうかを返す
func AllowedImageMIME(contentType string) (string, bool) {
	ext, ok := allowedImageMIMEs[contentType]
	return ext, ok
}

// ImageStorage はアップロード画像の保存先を抽象化する
type ImageStorage interface {
	Save(userID int32, contentType string, r io.Reader, size int64) (publicURL string, err error)
}

type LocalImageStorage struct {
	uploadDir string
	baseURL   string
}

func NewLocalImageStorage(uploadDir, baseURL string) (*LocalImageStorage, error) {
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return nil, fmt.Errorf("upload directory: %w", err)
	}
	return &LocalImageStorage{
		uploadDir: uploadDir,
		baseURL:   strings.TrimRight(baseURL, "/"),
	}, nil
}

func (s *LocalImageStorage) Save(userID int32, contentType string, r io.Reader, size int64) (string, error) {
	if size <= 0 {
		return "", fmt.Errorf("画像ファイルが空です")
	}
	if size > MaxImageSize {
		return "", fmt.Errorf("画像は5MB以下にしてください")
	}

	ext, ok := allowedImageMIMEs[contentType]
	if !ok {
		return "", fmt.Errorf("対応していない画像形式です（JPEG, PNG, GIF, WebP のみ）")
	}

	filename, err := randomFilename(userID, ext)
	if err != nil {
		return "", err
	}

	destPath := filepath.Join(s.uploadDir, filename)
	dest, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("画像の保存に失敗しました")
	}
	defer dest.Close()

	written, err := io.Copy(dest, io.LimitReader(r, MaxImageSize+1))
	if err != nil {
		_ = os.Remove(destPath)
		return "", fmt.Errorf("画像の保存に失敗しました")
	}
	if written > MaxImageSize {
		_ = os.Remove(destPath)
		return "", fmt.Errorf("画像は5MB以下にしてください")
	}

	return fmt.Sprintf("%s/uploads/%s", s.baseURL, filename), nil
}

func randomFilename(userID int32, ext string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("ファイル名の生成に失敗しました")
	}
	return fmt.Sprintf("%d_%s%s", userID, hex.EncodeToString(b), ext), nil
}
