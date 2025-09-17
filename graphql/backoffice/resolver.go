package self_service

import (
	"lavanilla/service/custom"
	"lavanilla/service/shopify"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CustomClient    *custom.Client
	ShopifyClient   *shopify.Client
	S3PresignClient *s3.PresignClient
}
