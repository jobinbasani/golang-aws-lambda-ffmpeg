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
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	lambda.Start(handler)
}

func handler(_ context.Context, s3Event events.S3Event) error {

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
	outputFile := strings.Replace(file.Name(), filepath.Ext(file.Name()), ".webm", 1)
	cmd := exec.Command("ffmpeg", "-i", file.Name(), outputFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	log.Printf("Execution output:\n%s\n", string(out))
	output, err := os.Open(outputFile)
	if err != nil {
		panic(err)
	}
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: &record.S3.Bucket.Name,
		Key:    aws.String(filepath.Base(outputFile)),
		Body:   output,
	})
	log.Printf("Copied %s to %s", outputFile, record.S3.Bucket.Name)
	return nil
}
