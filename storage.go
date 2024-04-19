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
	S3 StorageType = "S3"
)

func NewStorage(st StorageType) Storage {
	switch st {
	case S3:
		return NewS3Client()
	default:
		return nil
	}
}
