package volume

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
)

type VolumeService struct {
	S3Client *minio.Client
}

func NewService(s3 *minio.Client) VolumeService {
	return VolumeService{
		S3Client: s3,
	}
}

func (s VolumeService) CreateVolume(name string) error {
	return s.S3Client.MakeBucket(context.Background(), name, minio.MakeBucketOptions{})
}

func (s VolumeService) DeleteVolume(name string) error {
	return s.S3Client.RemoveBucket(context.Background(), name)
}

func (s VolumeService) GetUploadURL(bucket string, name string) (string, error) {
	url, err := s.S3Client.PresignedPutObject(context.Background(), bucket, name, 10*time.Minute)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s VolumeService) GetDownloadURL(bucket string, name string) (string, error) {
	url, err := s.S3Client.PresignedGetObject(context.Background(), bucket, name, 10*time.Minute, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
