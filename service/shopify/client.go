package shopify

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

const (
	shopId                   = "lvn-dev.myshopify.com"
	accessToken              = ""
	shopifyGraphqlApiVersion = "2025-10"
	headerKeyAccessToken     = "X-Shopify-Access-Token"
	headerKeyUserAgent       = "User-Agent"
	headerValueUserAgent     = "lavanilla/1.0"
)

type customHttpTransport struct {
	shopifyToken string
}

func (t *customHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(headerKeyUserAgent, headerValueUserAgent)
	req.Header.Add(headerKeyAccessToken, t.shopifyToken)
	return http.DefaultTransport.RoundTrip(req)
}

type Client struct {
	AccessToken string
	graphql     graphql.Client
}

func NewClient() *Client {
	endpoint := fmt.Sprintf("https://%v/admin/api/%v/graphql.json", shopId, shopifyGraphqlApiVersion)
	httpClient := &http.Client{Transport: &customHttpTransport{shopifyToken: accessToken}}
	return &Client{
		AccessToken: accessToken,
		graphql:     graphql.NewClient(endpoint, httpClient),
	}
}

func (c *Client) AddTag(ctx context.Context, orderId string, tag string) (*TagsAddResponse, error) {
	log.Printf("orderId: %v\n", orderId)
	return TagsAdd(ctx, c.graphql, orderId, tag)
}

func (c *Client) GetProductsSelfService(ctx context.Context) (*GetProductsSelfServiceResponse, error) {
	return GetProductsSelfService(ctx, c.graphql)
}

func (c *Client) GetProduct(ctx context.Context, id string) (*GetProductResponse, error) {
	return GetProduct(ctx, c.graphql, id)
}
