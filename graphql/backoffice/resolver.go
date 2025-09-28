package backoffice

import (
	"lavanilla/graphql/backoffice/controller/draft_order"
	"lavanilla/graphql/backoffice/controller/draft_order_product_variant"
	"lavanilla/service/custom"
	"lavanilla/service/shopify"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DraftOrder               *draft_order.Handler
	DraftOrderProductVariant *draft_order_product_variant.Handler
	CustomClient             *custom.Client
	ShopifyClient            *shopify.Client
	S3PresignClient          *s3.PresignClient
	S3Client                 *s3.Client
}
