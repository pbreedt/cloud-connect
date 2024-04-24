package aws

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

var bucketName string

func init() {
	bucketName = uuid.New().String()
}

func TestS3CreateBucket(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.CreateBucket(bucketName)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Bucket '%s' successfully created.\n", bucketName)
}

func TestS3Upload(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.StoreObject(bucketName, "test-object", "../test_data/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Test data successfully uploaded to bucket '%s'.\n", bucketName)
}

func TestS3ListBuckets(t *testing.T) {
	s3c := NewS3Client()
	t.Log("Buckets:")

	buckets, err := s3c.ListBuckets()
	if err != nil {
		t.Fatal(err)
	}
	for _, bucket := range buckets {
		t.Log(bucket)
	}
}

func TestS3ListBucketContent(t *testing.T) {
	s3c := NewS3Client()
	t.Logf("Objects in bucket '%s':\n", bucketName)

	objs, err := s3c.ListBucketContent(bucketName)
	if err != nil {
		t.Fatal(err)
	}
	for _, obj := range objs {
		t.Log(obj)
	}
}

func TestS3Download(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.RetrieveObject(bucketName, "test-object", "../test_data/aws_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test data successfully downloaded.")

	err = os.Remove("../test_data/aws_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test data successfully deleted from local storage.")
}

func TestS3Delete(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.DeleteObject(bucketName, []string{"test-object"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test data successfully deleted from bucket.")
}

func TestS3DeleteBucket(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.DeleteBucket(bucketName)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Bucket '%s' successfully deleted.\n", bucketName)
}
