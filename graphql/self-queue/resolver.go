package self_service

import "github.com/aws/aws-sdk-go-v2/service/s3"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	S3PresignClient *s3.PresignClient
	S3Client        *s3.Client
}
