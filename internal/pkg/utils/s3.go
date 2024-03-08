package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetS3PresignedURL(UID, bucketName string) (string, string, error) {
	const region = "eu-central-1"
	mySession := session.Must(session.NewSession())

	// Create a S3 client with additional configuration
	// TODO: pass region as a parameter
	svc := s3.New(mySession, aws.NewConfig().WithRegion(region))

	// Remove all spaces from the file name
	filePrefix := strings.ReplaceAll(UID, " ", "")

	// Generate a unique file name
	fileName := fmt.Sprintf("%s-%d", filePrefix, time.Now().Unix())

	// Generate a presigned URL for the file
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	// Generate a presigned URL for the file
	presignedURL, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "", "", err
	}

	storedFileURL := fmt.Sprintf("https:///%s.s3.%s.amazonaws.com/%s", bucketName, region, fileName)

	return presignedURL, storedFileURL, nil
}
