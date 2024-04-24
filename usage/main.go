package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/pbreedt/cloud-connect/storage"
)

func main() {
	bucketName := uuid.New().String()

	aws := storage.NewStorage(storage.Options{
		StorageType: storage.TypeS3,
	})
	log.Printf("Using AWS S3 storage provider\n")
	UseStorage(aws, bucketName)

	azure := storage.NewStorage(storage.Options{
		StorageType:          storage.TypeAzure,
		Azure_StorageAccount: "cs210032003763ea5a8",
	})
	log.Printf("\nUsing Azure storage provider\n")
	UseStorage(azure, bucketName)

	gcp := storage.NewStorage(storage.Options{
		StorageType:   storage.TypeGCP,
		GCP_ProjectId: "the-cloud-bootcamp-pfb",
	})
	log.Printf("\nUsing GCP storage provider\n")
	UseStorage(gcp, bucketName)
}

func UseStorage(cloudStorage storage.Storage, bucketName string) {
	err := cloudStorage.CreateBucket(bucketName)
	if err != nil {
		log.Fatalf("Error creating bucket(%s): %v\n", bucketName, err)
	}
	log.Printf("Successfully created bucket\n")

	buckets, err := cloudStorage.ListBuckets()
	if err != nil {
		log.Fatalf("Error listing buckets: %v\n", err)
	}
	for _, bucket := range buckets {
		if bucket == bucketName {
			log.Printf("List contains created bucket\n")
		}
	}

	err = cloudStorage.StoreObject(bucketName, "test-object-name", "../storage/test_data/testfile.txt")
	if err != nil {
		log.Fatalf("Error storing test object: %v\n", err)
	}
	log.Printf("Successfully stored test object")

	objects, err := cloudStorage.ListBucketContent(bucketName)
	if err != nil {
		log.Fatalf("Error listing bucket(%s) content: %v\n", bucketName, err)
	}
	for _, object := range objects {
		if object == "test-object-name" {
			log.Printf("List contains test object\n")
		}
	}

	err = cloudStorage.DeleteObject(bucketName, []string{"test-object-name"})
	if err != nil {
		log.Fatalf("Error deleting test object: %v\n", err)
	}
	log.Printf("Successfully deleted test object")

	err = cloudStorage.DeleteBucket(bucketName)
	if err != nil {
		log.Fatalf("Error deleting bucket(%s): %v\n", bucketName, err)
	}
	log.Printf("Successfully deleted bucket")

}
