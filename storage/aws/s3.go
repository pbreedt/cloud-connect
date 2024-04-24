package aws

import (
	"context"
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
	export AWS_DEFAULT_REGION=us-east-2
	us-east-1 not supported? see github.com/aws/aws-sdk-go-v2/service/s3/types.BucketLocationConstraint
*/

type S3Client struct {
	Client   *s3.Client
	location string
}

func NewS3Client() *S3Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return &S3Client{
		Client:   s3.NewFromConfig(cfg),
		location: cfg.Region,
	}
}

func (s3Client *S3Client) WithDefaultLocation(location string) *S3Client {
	s3Client.location = location
	return s3Client
}

func (s3Client *S3Client) CreateBucket(bucketName string) error {
	_, err := s3Client.Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(s3Client.location),
		},
	})
	return err
}

func (s3Client *S3Client) ListBuckets() ([]string, error) {
	buckets := []string{}

	result, err := s3Client.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return buckets, err
	}

	for _, bucket := range result.Buckets {
		buckets = append(buckets, *bucket.Name)
	}

	return buckets, nil
}

func (s3Client *S3Client) ListBucketContent(bucketName string) ([]string, error) {
	objects := []string{}

	objs, err := s3Client.Client.ListObjects(context.Background(), &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return objects, err
	}

	for _, object := range objs.Contents {
		objects = append(objects, *object.Key)
	}

	return objects, nil
}

func (s3Client *S3Client) DeleteBucket(bucketName string) error {
	_, err := s3Client.Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName)})
	return err
}

// For large files, use github.com/aws/aws-sdk-go-v2/feature/s3/manager.NewUploader
func (s3Client *S3Client) StoreObject(bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s3Client.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	return err
}

// For large files, use github.com/aws/aws-sdk-go-v2/feature/s3/manager.NewDownloader
// TODO: return []byte instead if writing to file
func (s3Client *S3Client) RetrieveObject(bucketName string, objectKey string, fileName string) error {
	result, err := s3Client.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return err
	}
	defer result.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}

	_, err = file.Write(body)
	return err
}

func (s3Client *S3Client) DeleteObject(bucketName string, objectKeys []string) error {
	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	_, err := s3Client.Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{Objects: objectIds},
	})

	return err
}
