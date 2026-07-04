package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func validateImage(size int64, contentType string) (ext string, err error) {
	if size <= 0 {
		return "", fmt.Errorf("画像ファイルが空です")
	}
	if size > MaxImageSize {
		return "", fmt.Errorf("画像は5MB以下にしてください")
	}

	ext, ok := allowedImageMIMEs[contentType]
	if !ok {
		return "", fmt.Errorf("対応していない画像形式です（JPEG, PNG のみ）")
	}
	return ext, nil
}

func readImageBody(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, MaxImageSize+1))
	if err != nil {
		return nil, fmt.Errorf("画像の読み込みに失敗しました")
	}
	if int64(len(data)) > MaxImageSize {
		return nil, fmt.Errorf("画像は5MB以下にしてください")
	}
	return data, nil
}

func objectKey(userID int32, ext string) (string, error) {
	filename, err := randomFilename(userID, ext)
	if err != nil {
		return "", err
	}
	return "uploads/" + filename, nil
}

func randomFilename(userID int32, ext string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("ファイル名の生成に失敗しました")
	}
	return fmt.Sprintf("%d_%s%s", userID, hex.EncodeToString(b), ext), nil
}
