package storage

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

var bcktName string

func init() {
	bcktName = uuid.New().String()
}

func TestS3CreateBucket(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.CreateBucket(bcktName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestS3Upload(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.StoreData(bcktName, "test-object", "./test_data/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestS3Download(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.RetrieveData(bcktName, "test-object", "./test_data/testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove("./test_data/testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestS3Delete(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.DeleteData(bcktName, []string{"test-object"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestS3DeleteBucket(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.DeleteBucket(bcktName)
	if err != nil {
		t.Fatal(err)
	}
}
