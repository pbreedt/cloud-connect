package storage

import (
	"github.com/pbreedt/cloud-connect/storage/aws"
	"github.com/pbreedt/cloud-connect/storage/azure"
	"github.com/pbreedt/cloud-connect/storage/gcp"
)

type Storage interface {
	CreateBucket(bucketName string) error
	ListBuckets() ([]string, error)
	ListBucketContent(bucketName string) ([]string, error)
	DeleteBucket(bucketName string) error

	StoreObject(bucketName string, objectKey string, fileName string) error
	RetrieveObject(bucketName string, objectKey string, fileName string) error
	DeleteObject(bucketName string, objectKeys []string) error
}

type StorageType string

const (
	TypeS3    StorageType = "S3"
	TypeGCP   StorageType = "GCP"
	TypeAzure StorageType = "AZURE"
)

type Options struct {
	StorageType StorageType
	ProjectId   string
	Location    string
}

func NewStorage(opts Options) Storage {
	switch opts.StorageType {
	case TypeS3:
		// AWS Location/Region retrieved from ~/.aws/config or env var AWS_DEFAULT_REGION
		return aws.NewS3Client()
	case TypeGCP:
		return gcp.NewCloudStorageClient(opts.ProjectId).WithDefaultLocation(opts.Location)
	case TypeAzure:
		return azure.NewBlobStorageClient(opts.ProjectId)
	default:
		return nil
	}
}
