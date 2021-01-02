package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, s3Event events.S3Event) error {

	record := s3Event.Records[0]

	key := record.S3.Object.Key
	sess, _ := session.NewSession(&aws.Config{Region: &record.AWSRegion})
	downloader := s3manager.NewDownloader(sess)
	file, err := os.Create(fmt.Sprintf("/tmp/%s", key))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: &record.S3.Bucket.Name,
			Key:    &key,
		})
	if err != nil {
		panic(err)
	}
	log.Printf("Downloaded %s", file.Name())
	return nil
}
