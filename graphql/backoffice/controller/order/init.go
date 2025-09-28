package order

import (
	"lavanilla/service/custom"
	"lavanilla/service/shopify"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Handler struct {
	shopifyClient   *shopify.Client
	customClient    *custom.Client
	s3Client        *s3.Client
	s3PresignClient *s3.PresignClient
}

func NewHandler(shopifyClient *shopify.Client, customClient *custom.Client, s3Client *s3.Client, s3PresignClient *s3.PresignClient) *Handler {
	return &Handler{
		shopifyClient:   shopifyClient,
		customClient:    customClient,
		s3Client:        s3Client,
		s3PresignClient: s3PresignClient,
	}
}
