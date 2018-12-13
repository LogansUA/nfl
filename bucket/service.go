package bucket

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/logansua/nfl_app/utils"
	"io"
	"mime/multipart"
	"os"
	"path"
)

type UploadFileToBucketResponse struct {
	Url string `json:"url"`
}

type Service interface {
	UploadPlayerAvatar(id uint, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

type service struct {
	Bucket *storage.BucketHandle
}

var (
	StorageBucketName string
	StorageBucket     *storage.BucketHandle
)

func New() Service {
	var err error

	StorageBucketName = os.Getenv("GOOGLE_CLOUD_BUCKET_NAME")

	if StorageBucketName == "" {
		panic(errors.New(fmt.Sprintf("Environement variable \"%s\" is not set!", "GOOGLE_CLOUD_BUCKET_NAME")))
	}

	StorageBucket, err = configureStorage(StorageBucketName)

	if err != nil {
		panic(err)
	}

	return &service{Bucket: StorageBucket}
}

func (s *service) UploadPlayerAvatar(id uint, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	name := fmt.Sprintf("%s%s", utils.RandToken(), path.Ext(fileHeader.Filename))
	filePath := fmt.Sprintf("%s/%d/%s", "players", id, name)

	_, err := uploadFileToBucket(file, fileHeader, filePath)

	if err != nil {
		return "", nil
	}

	return name, err
}

func uploadFileToBucket(file multipart.File, fileHeader *multipart.FileHeader, fullPath string) (url string, err error) {
	if StorageBucket == nil {
		return "", errors.New("storage bucket is missing - check config.go")
	}

	ctx := context.Background()
	writer := StorageBucket.Object(fullPath).NewWriter(ctx)

	// Warning: storage.AllUsers gives public read access to anyone.
	writer.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	writer.ContentType = fileHeader.Header.Get("Content-Type")

	// Entries are immutable, be aggressive about caching (1 day).
	writer.CacheControl = "public, max-age=86400"

	if _, err := io.Copy(writer, file); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", StorageBucketName, fullPath), nil
}

func configureStorage(bucketID string) (*storage.BucketHandle, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, err
	}

	return client.Bucket(bucketID), nil
}