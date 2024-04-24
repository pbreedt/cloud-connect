package gcp

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	bucketName string
	projectId  string
)

func init() {
	bucketName = uuid.New().String()
	projectId = "the-cloud-bootcamp-pfb"
}

func TestCSCreateBucket(t *testing.T) {
	gcp := NewCloudStorageClient(projectId)

	err := gcp.CreateBucket(bucketName)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Bucket '%s' successfully created.\n", bucketName)
}

func TestCSUpload(t *testing.T) {
	gcp := NewCloudStorageClient(projectId)

	err := gcp.StoreObject(bucketName, "test-object", "../test_data/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Test data successfully uploaded to bucket '%s'.\n", bucketName)
}

func TestCSListBuckets(t *testing.T) {
	gcp := NewCloudStorageClient(projectId)
	t.Logf("Buckets in project '%s':\n", projectId)

	buckets, err := gcp.ListBuckets()
	if err != nil {
		t.Fatal(err)
	}
	for _, bucket := range buckets {
		t.Log(bucket)
	}
}

func TestCSListBucketContent(t *testing.T) {
	gcp := NewCloudStorageClient(projectId)
	t.Logf("Objects in bucket '%s':\n", bucketName)

	objs, err := gcp.ListBucketContent(bucketName)
	if err != nil {
		t.Fatal(err)
	}
	for _, obj := range objs {
		t.Log(obj)
	}
}

func TestCSDownload(t *testing.T) {
	gcp := NewCloudStorageClient(projectId)

	err := gcp.RetrieveObject(bucketName, "test-object", "../test_data/gcp_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test data successfully downloaded.")

	err = os.Remove("../test_data/gcp_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test data successfully deleted from local storage.")
}

func TestCSDelete(t *testing.T) {
	gcp := NewCloudStorageClient(projectId)

	err := gcp.DeleteObject(bucketName, []string{"test-object"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test data successfully deleted from bucket.")
}

func TestCSDeleteBucket(t *testing.T) {
	gcp := NewCloudStorageClient(projectId)

	err := gcp.DeleteBucket(bucketName)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Bucket '%s' successfully deleted.\n", bucketName)
}
