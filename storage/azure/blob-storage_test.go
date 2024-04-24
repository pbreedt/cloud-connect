package azure

import (
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

func TestAZCreateBucket(t *testing.T) {
	az := NewBlobStorageClient(azStorageAccount)
	t.Logf("Blob storage client created: '%v'\n", az)
	err := az.CreateBucket(azBucketName)
	if err != nil {
		t.Fatal(err)
	}
}

// func TestGCPUpload(t *testing.T) {
// 	gcp := NewCloudStorageClient(gcpDefaultProjectId)

// 	err := gcp.StoreData(gcpBucketName, "test-object", "../test_data/testfile.txt")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestGCPDownload(t *testing.T) {
// 	gcp := NewCloudStorageClient(gcpDefaultProjectId)

// 	err := gcp.RetrieveData(gcpBucketName, "test-object", "../test_data/gcp_testfile_download.txt")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	err = os.Remove("../test_data/gcp_testfile_download.txt")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestGCPDelete(t *testing.T) {
// 	gcp := NewCloudStorageClient(gcpDefaultProjectId)

// 	err := gcp.DeleteData(gcpBucketName, []string{"test-object"})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

func TestAZDeleteBucket(t *testing.T) {
	az := NewBlobStorageClient(azStorageAccount)
	t.Logf("Blob storage client created: '%v'\n", az)
	err := az.DeleteBucket(azBucketName)
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
