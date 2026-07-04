package storage

import "os"

// NewImageStorageFromEnv は環境変数から S3 画像ストレージを生成する
func NewImageStorageFromEnv() (ImageStorage, error) {
	return NewS3ImageStorage(S3Config{
		Endpoint:   os.Getenv("S3_ENDPOINT"),
		Bucket:     os.Getenv("S3_BUCKET"),
		Region:     os.Getenv("S3_REGION"),
		AccessKey:  os.Getenv("S3_ACCESS_KEY"),
		SecretKey:  os.Getenv("S3_SECRET_KEY"),
		PublicBase: os.Getenv("S3_PUBLIC_BASE_URL"),
	})
}
