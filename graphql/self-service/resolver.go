package self_service

import (
	"lavanilla/service/custom"
	"lavanilla/service/shopify"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CustomClient  *custom.Client
	ShopifyClient *shopify.Client
}
