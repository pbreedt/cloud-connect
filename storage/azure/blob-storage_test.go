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

func TestCreateBucket(t *testing.T) {
	az := NewBlobStorageClient(azStorageAccount)
	t.Logf("Blob storage client created: '%v'\n", az)
	err := az.CreateBucket(azBucketName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpload(t *testing.T) {
	gcp := NewBlobStorageClient(azStorageAccount)

	err := gcp.StoreData(azBucketName, "test-object", "../test_data/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestListBucketContent(t *testing.T) {
	az := NewBlobStorageClient(azStorageAccount)
	t.Logf("Content of bucket: '%s'\n", azBucketName)
	az.ListBucketContents(azBucketName)
}

func TestDownload(t *testing.T) {
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

func TestDelete(t *testing.T) {
	gcp := NewBlobStorageClient(azStorageAccount)

	err := gcp.DeleteData(azBucketName, []string{"test-object"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteBucket(t *testing.T) {
	az := NewBlobStorageClient(azStorageAccount)
	t.Logf("Blob storage client created: '%v'\n", az)
	err := az.DeleteBucket(azBucketName)
	if err != nil {
		t.Fatal(err)
	}
}

// func TestListBuckets(t *testing.T) {
// 	az := NewBlobStorageClient(azStorageAccount)
// 	t.Logf("Buckets list :) %s\n", azBucketName)
// 	az.ListBuckets(azStorageAccount)
// }
