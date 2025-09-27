package draft_order_product_variant

import (
	"lavanilla/service/custom"
	"lavanilla/service/shopify"
)

type Handler struct {
	shopifyClient *shopify.Client
	customClient  *custom.Client
}

func NewHandler(shopifyClient *shopify.Client, customClient *custom.Client) *Handler {
	return &Handler{
		shopifyClient: shopifyClient,
		customClient:  customClient,
	}
}
