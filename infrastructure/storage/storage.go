package storage

import "io"

const MaxImageSize = 5 << 20 // 5MB

var allowedImageMIMEs = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
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
