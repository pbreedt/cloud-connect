package azure

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	azBucketName     string
	azStorageAccount string
)

func init() {
	azBucketName = uuid.New().String()
	azStorageAccount = "cs210032003763ea5a8"
}

func TestBSCreateBucket(t *testing.T) {
	az := NewBlobStorageClient(azStorageAccount)
	t.Logf("Blob storage client created: '%v'\n", az)
	err := az.CreateBucket(azBucketName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBSUpload(t *testing.T) {
	gcp := NewBlobStorageClient(azStorageAccount)

	err := gcp.StoreData(azBucketName, "test-object", "../test_data/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBSListBuckets(t *testing.T) {
	az := NewBlobStorageClient(azStorageAccount)
	t.Logf("Buckets list :) %s\n", azStorageAccount)
	az.ListBuckets()
}

func TestBSListBucketContent(t *testing.T) {
	az := NewBlobStorageClient(azStorageAccount)
	t.Logf("Content of bucket: '%s'\n", azBucketName)
	az.ListBucketContent(azBucketName)
}

func TestBSDownload(t *testing.T) {
	gcp := NewBlobStorageClient(azStorageAccount)

	err := gcp.RetrieveData(azBucketName, "test-object", "../test_data/az_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove("../test_data/az_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBSDelete(t *testing.T) {
	gcp := NewBlobStorageClient(azStorageAccount)

	err := gcp.DeleteData(azBucketName, []string{"test-object"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBSDeleteBucket(t *testing.T) {
	az := NewBlobStorageClient(azStorageAccount)
	t.Logf("Blob storage client created: '%v'\n", az)
	err := az.DeleteBucket(azBucketName)
	if err != nil {
		t.Fatal(err)
	}
}
