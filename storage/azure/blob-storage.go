package azure

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

/*
see:
	https://github.com/Azure/azure-sdk-for-go/
	https://learn.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-blobs-go

auth:
	https://github.com/azure-samples/azure-sdk-for-go-samples#prerequisites

	1. Create Subscription, get AZURE_SUBSCRIPTION_ID
	2. Create Storage Account
	3. Create Application Registration (Menu: "Microsoft Entra ID" > "App Registrations"), get AZURE_CLIENT_ID and AZURE_TENANT_ID
	4. Create Client Secret, get AZURE_CLIENT_SECRET

	AZURE_SUBSCRIPTION_ID=xxx
	AZURE_CLIENT_ID=xxx
	AZURE_TENANT_ID=xxx
	AZURE_CLIENT_SECRET=xxx
*/

type BlobStorageClient struct {
	Client         *azblob.Client
	storageAccount string
}

func NewBlobStorageClient(storageAccount string) *BlobStorageClient {

	url := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccount)

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}
	client, err := azblob.NewClient(url, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &BlobStorageClient{
		Client:         client,
		storageAccount: storageAccount,
	}

	// alternatively: use client factory
	// cred, err := azidentity.NewDefaultAzureCredential(nil)
	// if err != nil {
	// 	log.Fatalf("failed to obtain a credential: %v", err)
	// }
	// ctx := context.Background()
	// clientFactory, err := armstorage.NewClientFactory("<subscription-id>", cred, nil)
	// if err != nil {
	// 	log.Fatalf("failed to create client: %v", err)
	// }
	// res, err := clientFactory.NewBlobServicesClient()

	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func (az *BlobStorageClient) CreateBucket(bucketName string) error {
	_, err := az.Client.CreateContainer(context.Background(), bucketName, nil)
	return err
}

func (az *BlobStorageClient) ListBuckets() ([]string, error) {
	buckets := []string{}

	pager := az.Client.NewListContainersPager(&azblob.ListContainersOptions{})
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			return buckets, err
		}

		for _, container := range resp.ContainerItems {
			buckets = append(buckets, *container.Name)
		}
	}

	return buckets, nil
}

func (az *BlobStorageClient) ListBucketContent(bucketName string) ([]string, error) {
	objects := []string{}

	pager := az.Client.NewListBlobsFlatPager(bucketName, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Snapshots: true, Versions: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			return objects, err
		}

		for _, blob := range resp.Segment.BlobItems {
			objects = append(objects, *blob.Name)
		}
	}

	return objects, nil
}

func (az *BlobStorageClient) DeleteBucket(bucketName string) error {
	_, err := az.Client.DeleteContainer(context.TODO(), bucketName, nil)
	return err
}

func (az *BlobStorageClient) StoreObject(bucketName string, objectKey string, fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	_, err = az.Client.UploadBuffer(context.Background(), bucketName, objectKey, data, nil)
	return err
}

// TODO: return []byte instead if writing to file
func (az *BlobStorageClient) RetrieveObject(bucketName string, objectKey string, fileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	ds, err := az.Client.DownloadStream(ctx, bucketName, objectKey, nil)
	if err != nil {
		return err
	}

	retryReader := ds.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	err = retryReader.Close()
	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	body, err := io.ReadAll(retryReader)
	if err != nil {
		return err
	}

	_, err = file.Write(body)
	return err
}

func (az *BlobStorageClient) DeleteObject(bucketName string, objectKeys []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	for _, objectKey := range objectKeys {
		_, err := az.Client.DeleteBlob(ctx, bucketName, objectKey, nil)
		if err != nil {
			return fmt.Errorf("DeleteObject(%w)", err)
		}
	}

	return nil
}

/*
long running process:
ctx := context.Background()
// Call an asynchronous function to create a client. The return value is a poller object.
poller, err := client.BeginCreate(ctx, "resource_identifier", "additional_parameter")

if err != nil {
	// handle error...
}

// Call the poller object's PollUntilDone function that will block until the poller object
// has been updated to indicate the task has completed.
resp, err = poller.PollUntilDone(ctx, nil)
if err != nil {
	// handle error...
}

// Print the fact that the LRO completed.
fmt.Printf("LRO done")
*/
