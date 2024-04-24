package azure

import (
	"context"
	"fmt"
	"log"

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

// // For large files, use github.com/aws/aws-sdk-go-v2/feature/s3/manager.NewUploader
// func (az *BlobStorageClient) StoreData(bucketName string, objectKey string, fileName string) error {
// 	file, err := os.Open(fileName)
// 	if err != nil {
// 		log.Printf("Couldn't open file %v to upload. Error: %v\n", fileName, err)
// 	} else {
// 		defer file.Close()
// 		_, err = az.Client.PutObject(context.TODO(), &s3.PutObjectInput{
// 			Bucket: aws.String(bucketName),
// 			Key:    aws.String(objectKey),
// 			Body:   file,
// 		})
// 		if err != nil {
// 			log.Printf("Couldn't upload file %v to %v:%v. Error: %v\n",
// 				fileName, bucketName, objectKey, err)
// 		}
// 	}

// 	return err
// }

// // For large files, use github.com/aws/aws-sdk-go-v2/feature/s3/manager.NewDownloader
// // TODO: return []byte instead if writing to file
// func (az *BlobStorageClient) RetrieveData(bucketName string, objectKey string, fileName string) error {
// 	result, err := az.Client.GetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(objectKey),
// 	})
// 	if err != nil {
// 		log.Printf("Couldn't get object %v:%v. Error: %v\n", bucketName, objectKey, err)
// 		return err
// 	}
// 	defer result.Body.Close()
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		log.Printf("Couldn't create file %v. Error: %v\n", fileName, err)
// 		return err
// 	}
// 	defer file.Close()
// 	body, err := io.ReadAll(result.Body)
// 	if err != nil {
// 		log.Printf("Couldn't read object body from %v. Error: %v\n", objectKey, err)
// 	}
// 	_, err = file.Write(body)
// 	return err
// }

// func (az *BlobStorageClient) DeleteData(bucketName string, objectKeys []string) error {
// 	var objectIds []types.ObjectIdentifier
// 	for _, key := range objectKeys {
// 		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
// 	}
// 	_, err := az.Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
// 		Bucket: aws.String(bucketName),
// 		Delete: &types.Delete{Objects: objectIds},
// 	})
// 	if err != nil {
// 		log.Printf("Couldn't delete objects from bucket %v. Error: %v\n", bucketName, err)
// 		// } else {
// 		// 	log.Printf("Deleted %v objects.\n", len(output.Deleted))
// 	}
// 	return err
// }

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
