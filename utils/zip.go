package utils

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"lavanilla/service"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileData struct {
	Key  string
	Data []byte
	Err  error
}

func CreateZipArchive(ctx context.Context, zipFileName string, client *s3.Client, presignClient *s3.PresignClient, results <-chan FileData) (*string, error) {
	var zipBuf bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuf)

	for result := range results {
		if result.Err != nil {
			log.Printf("failed to download %s: %v\n", result.Key, result.Err)
			continue
		}

		// fw, err := zipWriter.Create(path.Base(result.Key))
		fw, err := zipWriter.Create(result.Key)
		if err != nil {
			log.Printf("failed to create zip entry: %v", err)
			continue
		}

		_, err = fw.Write(result.Data)
		if err != nil {
			log.Printf("failed to write %s to zip: %v\n", result.Key, err)
			continue
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %w", err)
	}

	zipKey := fmt.Sprintf("%s.zip", zipFileName)
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(service.S3BucketTemp),
		Key:    aws.String(zipKey),
		Body:   bytes.NewReader(zipBuf.Bytes()),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to upload zip to s3: %v", err))
	}

	object, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(service.S3BucketTemp),
		Key:    aws.String(zipKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to presign url: %v", err))
	}

	return &object.URL, nil
}
