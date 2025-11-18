package volume

import (
	"context"

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
