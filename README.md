# cloud-connect
> A package to wrap services from different cloud providers in a common, simplified interface

## General info
This main purpose of this package was to explore the different SDK's of AWS, GCP & Azure  

## Technologies
* Go - version 1.21.5

## Setup
Like with all Go modules, you can simply "go get" it.

```go get github.com/pbreedt/cloud-connect```

### AWS authentication
  1. Login to AWS console
  2. Go to or create users: IAM > Users
  3. Create and download Access key

	Store credentials in ~/.aws/credentials
	
  -- OR --  

  Set the following environment variables:
  ```
	export AWS_ACCESS_KEY_ID=xxx
	export AWS_SECRET_ACCESS_KEY=xxx
	export AWS_DEFAULT_REGION=us-east-2
  ```
	us-east-1 not supported? see github.com/aws/aws-sdk-go-v2/service/s3/types.BucketLocationConstraint

### GCP authentication
	1. Login to GCP console
	2. Create service account
	3. Create service account key (storing key in /path/to/sa-json.json file)
	4. export GOOGLE_APPLICATION_CREDENTIALS=/path/to/sa-json.json

### Azure authentication
	See: [env setup](https://github.com/azure-samples/azure-sdk-for-go-samples#prerequisites)  

	1. Create Subscription, get AZURE_SUBSCRIPTION_ID
	2. Create Storage Account
	3. Create Application Registration (Menu: "Microsoft Entra ID" > "App Registrations"), get AZURE_CLIENT_ID and AZURE_TENANT_ID
	4. Create Client Secret, get AZURE_CLIENT_SECRET

	AZURE_SUBSCRIPTION_ID=xxx
	AZURE_CLIENT_ID=xxx
	AZURE_TENANT_ID=xxx
	AZURE_CLIENT_SECRET=xxx

## Code Examples
Usage should be fairly straight forward:  
See [./usage/main.go](./usage/main.go)

```Go
package main

import (
	"fmt"

	"github.com/pbreedt/cloud-connect"
)

func main() {
	cloudStorage := storage.NewStorage(storage.Options{
        StorageType: storage.TypeS3,
    })

	err := cloudStorage.CreateBucket(bucketName)
}
```

## Features
Current features include:
* Creating bucket (or container in Azure)
* Listing buckets
* Listing bucket content
* Store object in bucket (from file source)
* Retrieve object from bucket (to file destination)
* Delete object from bucket
* Delete bucket

Possible improvements:
* Rather return io.Reader or []byte in RetrieveObject
* Rather provide io.Writer or []byte in StoreObject
* Implement:
  * Invoke lamda/cloud function
  * Access cloud db/table
  * Use cloud queueing
  * More

## Status
Project is: _in progress_

## Credits
Thanks to [ritaly](https://github.com/ritaly/README-cheatsheet) for a quick readme template

## Contact
Created by [@pbreedt](mailto:petrus.breedt@gmail.com)