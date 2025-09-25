package S3util

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func RenameS3Directory(ctx context.Context, s3client *s3.Client, bucket, oldPrefix, newPrefix string) error {

	// List objects under old prefix
	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(oldPrefix),
	}

	listPaginator := s3.NewListObjectsV2Paginator(s3client, listInput)

	for listPaginator.HasMorePages() {
		page, err := listPaginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, obj := range page.Contents {
			oldKey := *obj.Key
			newKey := strings.Replace(oldKey, oldPrefix, newPrefix, 1)

			// Copy object to new key
			_, err = s3client.CopyObject(ctx, &s3.CopyObjectInput{
				Bucket:     aws.String(bucket),
				CopySource: aws.String(bucket + "/" + oldKey),
				Key:        aws.String(newKey),
			})
			if err != nil {
				return fmt.Errorf("failed to copy %s to %s: %w", oldKey, newKey, err)
			}

			// Delete old object
			_, err = s3client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(oldKey),
			})
			if err != nil {
				return fmt.Errorf("failed to delete %s: %w", oldKey, err)
			}

			log.Printf("Moved %s → %s\n", oldKey, newKey)
		}
	}

	return nil
}
