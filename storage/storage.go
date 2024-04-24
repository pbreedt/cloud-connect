package storage

import (
	"log"

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
	StorageType          StorageType
	GCP_ProjectId        string
	Azure_StorageAccount string
	Location             string
}

func NewStorage(opts Options) Storage {
	if opts.StorageType == "" {
		log.Fatalf("Options.StorageType must be provided")
	}

	switch opts.StorageType {
	case TypeS3:
		// AWS Location/Region retrieved from ~/.aws/config or env var AWS_DEFAULT_REGION
		return aws.NewS3Client()
	case TypeGCP:
		if opts.GCP_ProjectId == "" {
			log.Fatalf("Options.GCP_ProjectId must be provided")
		}
		return gcp.NewCloudStorageClient(opts.GCP_ProjectId)
	case TypeAzure:
		if opts.Azure_StorageAccount == "" {
			log.Fatalf("Options.Azure_StorageAccount must be provided")
		}
		return azure.NewBlobStorageClient(opts.Azure_StorageAccount)
	default:
		return nil
	}
}
