package storage

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

func TestGCPCreateBucket(t *testing.T) {
	gcp := NewGCPClient(gcpDefaultProjectId)

	err := gcp.CreateBucket(gcpBucketName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGCPUpload(t *testing.T) {
	gcp := NewGCPClient(gcpDefaultProjectId)

	err := gcp.StoreData(gcpBucketName, "test-object", "./test_data/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGCPDownload(t *testing.T) {
	gcp := NewGCPClient(gcpDefaultProjectId)

	err := gcp.RetrieveData(gcpBucketName, "test-object", "./test_data/gcp_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove("./test_data/gcp_testfile_download.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGCPDelete(t *testing.T) {
	gcp := NewGCPClient(gcpDefaultProjectId)

	err := gcp.DeleteData(gcpBucketName, []string{"test-object"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGCPDeleteBucket(t *testing.T) {
	gcp := NewGCPClient(gcpDefaultProjectId)

	err := gcp.DeleteBucket(gcpBucketName)
	if err != nil {
		t.Fatal(err)
	}
}

// func TestGCPList(t *testing.T) {
// 	gcp := NewGCPClient()
// 	t.Logf("Bucket name: '%v'\n", gcpBucketName)
// 	gcp.ListBucketContents(gcpBucketName)
// }

// func TestGCPListBuckets(t *testing.T) {
// 	gcp := NewGCPClient()
// 	t.Logf("Bucket name: '%v'\n", gcpBucketName)
// 	gcp.ListBuckets("the-cloud-bootcamp-pfb")
// }
