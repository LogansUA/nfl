package file

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/GoogleCloudPlatform/golang-samples/getting-started/bookshelf"
	"github.com/logansua/nfl_app/utils"
	"io"
	"mime/multipart"
	"os"
	"path"
)

type UploadFileToBucketResponse struct {
	Url string `json:"url"`
}

var (
	StorageBucketName string
	StorageBucket     *storage.BucketHandle
)

func ConfigureBucketStorage() error {
	var err error

	StorageBucketName = os.Getenv("GOOGLE_CLOUD_BUCKET_NAME")

	if StorageBucketName == "" {
		return errors.New(fmt.Sprintf("Environement variable \"%s\" is not set!", "GOOGLE_CLOUD_BUCKET_NAME"))
	}

	StorageBucket, err = configureStorage(StorageBucketName)

	return err
}

func UploadFileToBucket(file multipart.File, fileHeader *multipart.FileHeader) (url string, err error) {
	if StorageBucket == nil {
		return "", errors.New("storage bucket is missing - check config.go")
	}

	// random filename, retaining existing extension.
	name := utils.RandToken() + path.Ext(fileHeader.Filename)

	ctx := context.Background()
	writer := StorageBucket.Object(name).NewWriter(ctx)

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

	const publicURL = "https://storage.googleapis.com/%s/%s"

	return fmt.Sprintf(publicURL, bookshelf.StorageBucketName, name), nil
}

func configureStorage(bucketID string) (*storage.BucketHandle, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketID), nil
}
