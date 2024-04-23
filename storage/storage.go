package storage

type Storage interface {
	CreateBucket(bucketName string) error
	DeleteBucket(bucketName string) error
	StoreData(bucketName string, objectKey string, fileName string) error
	RetrieveData(bucketName string, objectKey string, fileName string) error
	DeleteData(bucketName string, objectKeys []string) error
}

type StorageType string

const (
	TypeS3  StorageType = "S3"
	TypeGCP StorageType = "GCP"
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
		return NewS3Client()
	case TypeGCP:
		return NewGCPClient(opts.ProjectId).WithDefaultLocation(opts.Location)
	default:
		return nil
	}
}
