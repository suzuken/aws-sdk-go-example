package main

import (
	"flag"
	"fmt"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awserr"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/service/s3"
	"log"
	"os"
)

var BucketName = flag.String("bucket", "your-example-bucket-name", "your example bucket name")
var region = flag.String("region", "ap-northeast-1", "s3 region")
var key = flag.String("key", "your-example-s3-key-name", "key of s3 path")

// example script using s3
// put README.md, list, head, copy and delete
//
// Usage: ./s3_misc --bucket path-to-your-bucket --key example-key
func main() {
	flag.Parse()
	// credentials are automatically loaded from ~/.aws/credentials
	s := s3.New(&aws.Config{Region: *region})

	f, err := os.Open("./README.md")
	if err != nil {
		log.Fatalf("[err] %s", err)
	}
	defer f.Close()

	// http://godoc.org/github.com/awslabs/aws-sdk-go/service/s3#example-S3-PutObject
	params := &s3.PutObjectInput{
		Body: f,
		// http://docs.aws.amazon.com/AmazonS3/latest/dev/UsingMetadata.html
		Key:    aws.String(*key + "/path-to-readme.md"),
		Bucket: aws.String(*BucketName),
	}
	resp, err := s.PutObject(params)
	handleError(err)
	fmt.Println(awsutil.StringValue(resp))
	// {
	//     ETag: "\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\""
	// }

	// http://godoc.org/github.com/awslabs/aws-sdk-go/service/s3#example-S3-ListObjects
	listParams := &s3.ListObjectsInput{
		Bucket: aws.String(*BucketName),
		Prefix: aws.String(*key),
	}
	listResp, err := s.ListObjects(listParams)
	handleError(err)
	fmt.Println(awsutil.StringValue(listResp))
	// {
	//   Contents: [{
	//       ETag: "\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\"",
	//       Key: "example-key/path-to-readme.md",
	//       LastModified: 2015-05-18 10:59:14 +0000 UTC,
	//       Owner: {
	//         DisplayName: "hogehoge",
	//         ID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	//       },
	//       Size: 37,
	//       StorageClass: "STANDARD"
	//     }],
	//   IsTruncated: false,
	//   Marker: "",
	//   MaxKeys: 1000,
	//   Name: "example-bucket",
	//   Prefix: "example-key"
	// }

	headParams := &s3.HeadObjectInput{
		Bucket: aws.String(*BucketName),
		Key:    aws.String(*key + "/path-to-readme.md"),
	}
	headResp, err := s.HeadObject(headParams)
	handleError(err)
	fmt.Println(awsutil.StringValue(headResp))
	// {
	// AcceptRanges: "bytes",
	// ContentLength: 37,
	// ContentType: "binary/octet-stream",
	// ETag: "\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\"",
	// LastModified: 2015-05-21 07:53:20 +0000 UTC,
	// Metadata: {
	//
	// }

	copyParams := &s3.CopyObjectInput{
		Bucket:     aws.String(*BucketName),
		CopySource: aws.String(fmt.Sprintf("%s/%s/path-to-readme.md", *BucketName, *key)),
		Key:        aws.String(*key + "/path-to-readme-copy.md"),
	}
	copyResp, err := s.CopyObject(copyParams)
	handleError(err)
	fmt.Println(awsutil.StringValue(copyResp))
	// {
	  // CopyObjectResult: {
		// ETag: "\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\"",
		// LastModified: 2015-05-21 07:55:07 +0000 UTC
	  // }
	// }

	deleteParams := &s3.DeleteObjectInput{
		Bucket:     aws.String(*BucketName),
		Key:        aws.String(*key + "/path-to-readme-copy.md"),
	}
	deleteResp, err := s.DeleteObject(deleteParams)
	handleError(err)
	fmt.Println(awsutil.StringValue(deleteResp))
}

// error handling helper
func handleError(err error) {
	// match type of error: https://github.com/awslabs/aws-sdk-go/issues/238
	if awsErr, ok := err.(awserr.Error); ok {
		log.Println(awsErr.Code(), awsErr.Message(), awsErr.Error(), awsErr.OrigErr())
		if reqErr, ok := err.(awserr.RequestFailure); ok {
			log.Println(reqErr.StatusCode(), reqErr.RequestID())
		}
	} else if err != nil {
		panic(err)
	}
}
