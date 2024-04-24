package azure

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	bucketName     string
	storageAccount string
)

func init() {
	bucketName = uuid.New().String()
	storageAccount = "cs210032003763ea5a8"
}

func TestBSCreateBucket(t *testing.T) {
	az := NewBlobStorageClient(storageAccount)
	t.Logf("Blob storage client created: '%v'\n", az)
	err := az.CreateBucket(bucketName)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Bucket '%s' successfully created.\n", bucketName)
}

func TestBSUpload(t *testing.T) {
	gcp := NewBlobStorageClient(storageAccount)

	err := gcp.StoreObject(bucketName, "test-object", "../test_data/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Test data successfully uploaded to bucket '%s'.\n", bucketName)
}

func TestBSListBuckets(t *testing.T) {
	az := NewBlobStorageClient(storageAccount)
	t.Logf("Buckets in storage account '%s':\n", storageAccount)

	buckets, err := az.ListBuckets()
	if err != nil {
		t.Fatal(err)
	}
	for bkt := range buckets {
		t.Log(bkt)
	}
}

func TestBSListBucketContent(t *testing.T) {
	az := NewBlobStorageClient(storageAccount)
	t.Logf("Objects in bucket '%s':\n", bucketName)

	objs, err := az.ListBucketContent(bucketName)
	if err != nil {
		t.Fatal(err)
	}
	for obj := range objs {
		t.Log(obj)
	}
}

func TestBSDownload(t *testing.T) {
	gcp := NewBlobStorageClient(storageAccount)

	err := gcp.RetrieveObject(bucketName, "test-object", "../test_data/az_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test data successfully downloaded.")

	err = os.Remove("../test_data/az_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test data successfully deleted from local storage.")
}

func TestBSDelete(t *testing.T) {
	gcp := NewBlobStorageClient(storageAccount)

	err := gcp.DeleteObject(bucketName, []string{"test-object"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test data successfully deleted from bucket.")
}

func TestBSDeleteBucket(t *testing.T) {
	az := NewBlobStorageClient(storageAccount)

	err := az.DeleteBucket(bucketName)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Bucket '%s' successfully deleted.\n", bucketName)
}
