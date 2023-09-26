package utils

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

func UploadVideoToS3(videoData []byte, s3Path string) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	// Create a new AWS session with the loaded access keys
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
		// Add more configurations if needed.
	})
	if err != nil {
		return err
	}

	// Create an S3 uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the video data to S3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("reanweb"),
		Key:    aws.String(s3Path),
		Body:   bytes.NewReader(videoData),
	})
	if err != nil {
		return err
	}

	return nil
}
