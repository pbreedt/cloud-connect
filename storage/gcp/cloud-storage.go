package gcp

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

/*
see:
	https://cloud.google.com/docs
	https://cloud.google.com/storage/docs/quickstart

auth:
	1. Login to GCP console
	2. Create service account
	3. Create service account key (storing key in /path/to/sa-json.json file)
	4. export GOOGLE_APPLICATION_CREDENTIALS=/path/to/sa-json.json
*/

type CloudStorageClient struct {
	Client    *storage.Client
	projectId string
	location  string
}

func NewCloudStorageClient(projectId string) *CloudStorageClient {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return &CloudStorageClient{
		Client:    client,
		projectId: projectId,
		location:  "US-CENTRAL1",
	}
}

func (gcpClient *CloudStorageClient) WithDefaultLocation(location string) *CloudStorageClient {
	gcpClient.location = location
	return gcpClient
}

// ################
// Bucket functions
// ################
func (gcpClient *CloudStorageClient) CreateBucket(bucketName string) error {
	bkt := gcpClient.Client.Bucket(bucketName)
	attrLocation := &storage.BucketAttrs{
		Location: gcpClient.location,
	}

	err := bkt.Create(context.Background(), gcpClient.projectId, attrLocation)
	if err != nil {
		return err
	}

	return nil
}

func (gcpClient *CloudStorageClient) ListBuckets() ([]string, error) {
	buckets := []string{}

	it := gcpClient.Client.Buckets(context.Background(), gcpClient.projectId)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return buckets, err
		}
		buckets = append(buckets, attrs.Name)
	}

	return buckets, nil
}

func (gcpClient *CloudStorageClient) ListBucketContent(bucketName string) ([]string, error) {
	objects := []string{}

	bkt := gcpClient.Client.Bucket(bucketName)
	it := bkt.Objects(context.Background(), nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return objects, err
		}
		objects = append(objects, attrs.Name)
	}

	return objects, nil
}

func (gcpClient *CloudStorageClient) DeleteBucket(bucketName string) error {
	bkt := gcpClient.Client.Bucket(bucketName)

	err := bkt.Delete(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// ########################
// Bucket content functions
// ########################

// uploadFile uploads an object.
func (gcpClient *CloudStorageClient) StoreObject(bucketName string, objectKey string, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	o := gcpClient.Client.Bucket(bucketName).Object(objectKey)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	// o = o.If(storage.Conditions{DoesNotExist: true})

	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %w", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

// downloadFile downloads an object to a file.
func (gcpClient *CloudStorageClient) RetrieveObject(bucketName string, objectKey string, fileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	rc, err := gcpClient.Client.Bucket(bucketName).Object(objectKey).NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	return nil

}

func (gcpClient *CloudStorageClient) DeleteObject(bucketName string, objectKeys []string) error {
	for _, objectKey := range objectKeys {
		err := gcpClient.Client.Bucket(bucketName).Object(objectKey).Delete(context.Background())
		if err != nil {
			return fmt.Errorf("DeleteObject(%w)", err)
		}
	}
	return nil
}
