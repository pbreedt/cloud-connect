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
	// defaultLocation string
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
		// defaultLocation: "US-CENTRAL1",
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
	if err != nil {
		log.Printf("Couldn't create bucket %v in storage account %s. Error: %v\n",
			bucketName, az.storageAccount, err)
	}
	return err
}

func (az *BlobStorageClient) DeleteBucket(bucketName string) error {
	_, err := az.Client.DeleteContainer(context.TODO(), bucketName, nil)
	if err != nil {
		log.Printf("Couldn't delete bucket %v. Error: %v\n", bucketName, err)
	}
	return err
}

// For large files, use github.com/aws/aws-sdk-go-v2/feature/s3/manager.NewUploader
func (az *BlobStorageClient) StoreData(bucketName string, objectKey string, fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Error: %v\n", fileName, err)
	} else {
		_, err = az.Client.UploadBuffer(context.Background(), bucketName, objectKey, data, nil) // &azblob.UploadBufferOptions{}
		if err != nil {
			log.Printf("Couldn't upload file %v to %v:%v. Error: %v\n",
				fileName, bucketName, objectKey, err)
		}
	}

	return err
}

// TODO: return []byte instead if writing to file
func (az *BlobStorageClient) RetrieveData(bucketName string, objectKey string, fileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	ds, err := az.Client.DownloadStream(ctx, bucketName, objectKey, nil)
	if err != nil {
		return fmt.Errorf("DownloadStream(%s): %w", objectKey, err)
	}

	retryReader := ds.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	err = retryReader.Close()
	if err != nil {
		return fmt.Errorf("DownloadStream(%s).NewRetryReader(): %w", objectKey, err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Error: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(retryReader)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Error: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}

func (az *BlobStorageClient) DeleteData(bucketName string, objectKeys []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()
	for _, objectKey := range objectKeys {
		_, err := az.Client.DeleteBlob(ctx, bucketName, objectKey, nil)
		if err != nil {
			return fmt.Errorf("DeleteBlob(%s): %w", objectKey, err)
		}
		log.Printf("Data deleted from %s/%s\n", bucketName, objectKey)
	}
	return nil
}

func (az *BlobStorageClient) ListBucketContents(bucketName string) {
	// List the blobs in the container
	pager := az.Client.NewListBlobsFlatPager(bucketName, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Snapshots: true, Versions: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		for _, blob := range resp.Segment.BlobItems {
			fmt.Println(*blob.Name)
		}
	}
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
