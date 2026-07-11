package storage

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3ImageStorage は S3 互換ストレージ（開発時は S3Mock）に画像を保存する
type S3ImageStorage struct {
	client    *s3.S3
	bucket    string
	publicURL string
}

// S3Config は S3 クライアントの設定
type S3Config struct {
	Endpoint   string
	Bucket     string
	Region     string
	AccessKey  string
	SecretKey  string
	PublicBase string // ブラウザ向け URL（例: http://localhost:9090）
}

func NewS3ImageStorage(cfg S3Config) (*S3ImageStorage, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("S3_ENDPOINT が設定されていません")
	}
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("S3_BUCKET が設定されていません")
	}
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	if cfg.AccessKey == "" {
		cfg.AccessKey = "test"
	}
	if cfg.SecretKey == "" {
		cfg.SecretKey = "test"
	}
	if cfg.PublicBase == "" {
		cfg.PublicBase = strings.TrimRight(cfg.Endpoint, "/")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.Region),
		Endpoint:         aws.String(cfg.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		),
		DisableSSL: aws.Bool(strings.HasPrefix(cfg.Endpoint, "http://")),
	})
	if err != nil {
		return nil, fmt.Errorf("S3 セッション: %w", err)
	}

	client := s3.New(sess)
	storage := &S3ImageStorage{
		client:    client,
		bucket:    cfg.Bucket,
		publicURL: strings.TrimRight(cfg.PublicBase, "/"),
	}

	if err := storage.ensureBucket(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *S3ImageStorage) ensureBucket() error {
	_, err := s.client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err == nil {
		return nil
	}

	if aerr, ok := err.(awserr.Error); ok && (aerr.Code() == s3.ErrCodeNoSuchBucket || aerr.Code() == "NotFound") {
		_, err = s.client.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(s.bucket),
		})
		if err != nil {
			return fmt.Errorf("S3 バケット作成: %w", err)
		}
		return nil
	}

	return fmt.Errorf("S3 バケット確認: %w", err)
}

func (s *S3ImageStorage) Save(userID int32, contentType string, r io.Reader, size int64) (string, error) {
	ext, err := validateImage(size, contentType)
	if err != nil {
		return "", err
	}

	key, err := objectKey(userID, ext)
	if err != nil {
		return "", err
	}

	body, err := readImageBody(r)
	if err != nil {
		return "", err
	}

	_, err = s.client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(body),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(len(body))),
	})
	if err != nil {
		return "", fmt.Errorf("画像の保存に失敗しました")
	}

	return s.buildPublicURL(key), nil
}

func (s *S3ImageStorage) PresignUpload(userID int32, contentType string, size int64) (*PresignUploadResult, error) {
	ext, err := validateImage(size, contentType)
	if err != nil {
		return nil, err
	}

	key, err := objectKey(userID, ext)
	if err != nil {
		return nil, err
	}

	req, _ := s.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	})

	uploadURL, err := req.Presign(15 * time.Minute)
	if err != nil {
		return nil, fmt.Errorf("アップロードURLの発行に失敗しました")
	}

	return &PresignUploadResult{
		Key:       key,
		UploadURL: uploadURL,
		PublicURL: s.buildPublicURL(key),
	}, nil
}

func (s *S3ImageStorage) ConfirmUpload(key string) (string, error) {
	head, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && (aerr.Code() == s3.ErrCodeNoSuchKey || aerr.Code() == "NotFound") {
			return "", fmt.Errorf("画像のアップロードが確認できませんでした")
		}
		return "", fmt.Errorf("画像の確認に失敗しました")
	}

	if head.ContentLength == nil || *head.ContentLength <= 0 {
		return "", fmt.Errorf("画像ファイルが空です")
	}
	if *head.ContentLength > MaxImageSize {
		return "", fmt.Errorf("画像は5MB以下にしてください")
	}

	contentType := ""
	if head.ContentType != nil {
		contentType = *head.ContentType
	}
	if _, ok := allowedImageMIMEs[contentType]; !ok {
		return "", fmt.Errorf("対応していない画像形式です（JPEG, PNG のみ）")
	}

	return s.buildPublicURL(key), nil
}

func (s *S3ImageStorage) buildPublicURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucket, key)
}
