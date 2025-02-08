package repository

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Repository struct {
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	bucketName string
}

func NewS3Repository(bucketName string) *S3Repository {
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	token := os.Getenv("AWS_SESSION_TOKEN")

	if awsAccessKeyID == "" || awsSecretAccessKey == "" || awsRegion == "" || token == "" {
		log.Fatal("Missing required AWS environment variables")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, token),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	return &S3Repository{
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
		bucketName: bucketName,
	}
}

func (r *S3Repository) UploadFile(key string, file io.Reader) error {
	_, err := r.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	return err
}

func (r *S3Repository) DownloadFile(key string) (io.Reader, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := r.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Println("Error downloading file", err)
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
