package gcp

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	gcpBucketName       string
	gcpDefaultProjectId string
)

func init() {
	gcpBucketName = uuid.New().String()
	gcpDefaultProjectId = "the-cloud-bootcamp-pfb"
}

func TestCSCreateBucket(t *testing.T) {
	gcp := NewCloudStorageClient(gcpDefaultProjectId)

	err := gcp.CreateBucket(gcpBucketName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCSUpload(t *testing.T) {
	gcp := NewCloudStorageClient(gcpDefaultProjectId)

	err := gcp.StoreData(gcpBucketName, "test-object", "../test_data/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCSListBuckets(t *testing.T) {
	gcp := NewCloudStorageClient(gcpDefaultProjectId)
	t.Logf("Buckets list :) %s\n", gcpDefaultProjectId)
	gcp.ListBuckets()
}

func TestCSListBucketContent(t *testing.T) {
	gcp := NewCloudStorageClient(gcpDefaultProjectId)
	t.Logf("Content of bucket: '%s'\n", gcpBucketName)
	gcp.ListBucketContent(gcpBucketName)
}

func TestCSDownload(t *testing.T) {
	gcp := NewCloudStorageClient(gcpDefaultProjectId)

	err := gcp.RetrieveData(gcpBucketName, "test-object", "../test_data/gcp_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove("../test_data/gcp_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCSDelete(t *testing.T) {
	gcp := NewCloudStorageClient(gcpDefaultProjectId)

	err := gcp.DeleteData(gcpBucketName, []string{"test-object"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCSDeleteBucket(t *testing.T) {
	gcp := NewCloudStorageClient(gcpDefaultProjectId)

	err := gcp.DeleteBucket(gcpBucketName)
	if err != nil {
		t.Fatal(err)
	}
}
