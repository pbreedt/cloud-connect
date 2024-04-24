package aws

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

/*
see:

	https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/gov2
	https://aws.github.io/aws-sdk-go-v2/docs/getting-started/

auth:

	have credentials in ~/.aws/credentials
	-- or --
	export AWS_ACCESS_KEY_ID=xxx
	export AWS_SECRET_ACCESS_KEY=xxx
	export AWS_DEFAULT_REGION=us-east-1
*/
type S3Client struct {
	Client *s3.Client
}

func NewS3Client() *S3Client {
	// config.NewEnvConfig()
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	return &S3Client{Client: s3.NewFromConfig(cfg)}
}

func (s3Client *S3Client) CreateBucket(bucketName string) error {
	_, err := s3Client.Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		// CreateBucketConfiguration: &types.CreateBucketConfiguration{
		// 	LocationConstraint: types.BucketLocationConstraint(region),
		// },
	})
	if err != nil {
		log.Printf("Couldn't create bucket %v in Region %v. Error: %v\n",
			bucketName, s3Client.Client.Options().Region, err)
	}
	return err
}

func (s3Client *S3Client) ListBuckets() {
	result, err := s3Client.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		fmt.Printf("Couldn't list buckets for your account. Error: %v\n", err)
		return
	}
	if len(result.Buckets) == 0 {
		fmt.Println("You don't have any buckets!")
	} else {
		for _, bucket := range result.Buckets {
			fmt.Printf("\t%v\n", *bucket.Name)
		}
	}
}

func (s3Client *S3Client) ListBucketContent(bucketName string) {
	objects, err := s3Client.Client.ListObjects(context.Background(), &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Found %v objects.\n", len(objects.Contents))
	// var objKeys []string
	for _, object := range objects.Contents {
		// objKeys = append(objKeys, *object.Key)
		log.Printf("\t%v\n", *object.Key)
	}
}

func (s3Client *S3Client) DeleteBucket(bucketName string) error {
	_, err := s3Client.Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName)})
	if err != nil {
		log.Printf("Couldn't delete bucket %v. Error: %v\n", bucketName, err)
	}
	return err
}

// For large files, use github.com/aws/aws-sdk-go-v2/feature/s3/manager.NewUploader
func (s3Client *S3Client) StoreData(bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Error: %v\n", fileName, err)
	} else {
		defer file.Close()
		_, err = s3Client.Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			log.Printf("Couldn't upload file %v to %v:%v. Error: %v\n",
				fileName, bucketName, objectKey, err)
		}
	}

	return err
}

// For large files, use github.com/aws/aws-sdk-go-v2/feature/s3/manager.NewDownloader
// TODO: return []byte instead if writing to file
func (s3Client *S3Client) RetrieveData(bucketName string, objectKey string, fileName string) error {
	result, err := s3Client.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Error: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Error: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Error: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}

func (s3Client *S3Client) DeleteData(bucketName string, objectKeys []string) error {
	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	_, err := s3Client.Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		log.Printf("Couldn't delete objects from bucket %v. Error: %v\n", bucketName, err)
		// } else {
		// 	log.Printf("Deleted %v objects.\n", len(output.Deleted))
	}
	return err
}
