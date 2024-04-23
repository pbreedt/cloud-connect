package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

/*
	1. Login to GCP console
	2. Create service account
	3. Create service account key (storing key in /path/to/sa-json.json file)
	4. export GOOGLE_APPLICATION_CREDENTIALS=/path/to/sa-json.json
*/

type GCPClient struct {
	Client          *storage.Client
	projectId       string
	defaultLocation string
}

func NewGCPClient(defaultProjectId string) *GCPClient {
	client, err := storage.NewClient(context.Background())

	if err != nil {
		log.Fatal(err)
	}
	return &GCPClient{
		Client:          client,
		projectId:       defaultProjectId,
		defaultLocation: "US-CENTRAL1",
	}
}

func (gcpClient *GCPClient) WithDefaultLocation(location string) *GCPClient {
	gcpClient.defaultLocation = location
	return gcpClient
}

func (gcpClient *GCPClient) CreateBucket(bucketName string) error {
	bkt := gcpClient.Client.Bucket(bucketName)
	attrLocation := &storage.BucketAttrs{
		Location: gcpClient.defaultLocation,
	}
	err := bkt.Create(context.Background(), gcpClient.projectId, attrLocation)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("Bucket %v created\n", bucketName)
	return nil
}

func (gcpClient *GCPClient) DeleteBucket(bucketName string) error {
	bkt := gcpClient.Client.Bucket(bucketName)
	err := bkt.Delete(context.Background())
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("Bucket %v deleted\n", bucketName)
	return nil
}

// uploadFile uploads an object.
func (gcpClient *GCPClient) StoreData(bucketName string, objectKey string, fileName string) error {
	// Open local file.
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
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
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	log.Printf("Data stored in %s/%s\n", bucketName, objectKey)
	return nil
}

// downloadFile downloads an object to a file.
func (gcpClient *GCPClient) RetrieveData(bucketName string, objectKey string, fileName string) error {
	ctx := context.Background()
	// client, err := storage.NewClient(ctx)
	// if err != nil {
	// 	return fmt.Errorf("storage.NewClient: %w", err)
	// }
	// defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}

	rc, err := gcpClient.Client.Bucket(bucketName).Object(objectKey).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %w", objectKey, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("f.Close: %w", err)
	}

	fmt.Printf("Data from %s/%s downloaded to local file %s\n", bucketName, objectKey, fileName)

	return nil

}

func (gcpClient *GCPClient) DeleteData(bucketName string, objectKeys []string) error {
	for _, objectKey := range objectKeys {
		err := gcpClient.Client.Bucket(bucketName).Object(objectKey).Delete(context.Background())
		if err != nil {
			return fmt.Errorf("Object(%q).Delete: %w", objectKey, err)
		}
		log.Printf("Data deleted from %s/%s\n", bucketName, objectKey)
	}
	return nil
}

// func (gcpClient *GCPClient) ListBucketContents(bucketName string) {
// 	bkt := gcpClient.Client.Bucket(bucketName)
// 	it := bkt.Objects(context.Background(), nil)
// 	for {
// 		attrs, err := it.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(attrs.Name, attrs.ContentType)
// 	}
// }

// func (gcpClient *GCPClient) ListBuckets(prjId string) {
// 	it := gcpClient.Client.Buckets(context.Background(), prjId)
// 	for {
// 		attrs, err := it.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(attrs.Name, attrs.Location)
// 	}
// }
