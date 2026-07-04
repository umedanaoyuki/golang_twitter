package storage

import (
	"bytes"
	"fmt"
	"io"
	"strings"

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

	return fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucket, key), nil
}
