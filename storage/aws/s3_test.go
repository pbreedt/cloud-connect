package aws

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

var awsBucketName string

func init() {
	awsBucketName = uuid.New().String()
}

func TestS3CreateBucket(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.CreateBucket(awsBucketName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestS3Upload(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.StoreData(awsBucketName, "test-object", "../test_data/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestS3ListBuckets(t *testing.T) {
	s3c := NewS3Client()
	t.Logf("Buckets list :) %s\n", s3c.Client.Options().AppID)
	s3c.ListBuckets()
}

func TestS3ListBucketContent(t *testing.T) {
	s3c := NewS3Client()
	t.Logf("Content of bucket: '%s'\n", awsBucketName)
	s3c.ListBucketContent(awsBucketName)
}

func TestS3Download(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.RetrieveData(awsBucketName, "test-object", "../test_data/aws_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove("../test_data/aws_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestS3Delete(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.DeleteData(awsBucketName, []string{"test-object"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestS3DeleteBucket(t *testing.T) {
	s3c := NewS3Client()

	err := s3c.DeleteBucket(awsBucketName)
	if err != nil {
		t.Fatal(err)
	}
}
