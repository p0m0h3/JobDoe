package volume

import (
	"context"
	"time"
)

type DockerVolumeService struct {
	DockerClient *docker.Client
}

func (s DockerVolumeService) GetUploadURL(bucket string, name string) (string, error) {
	url, err := s.DockerClient.PresignedPutObject(context.Background(), bucket, name, 10*time.Minute)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s DockerVolumeService) GetDownloadURL(bucket string, name string) (string, error) {
	url, err := s.DockerClient.PresignedGetObject(context.Background(), bucket, name, 10*time.Minute, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
