package shopify

import (
	"context"
	"fmt"
	"lavanilla/service"
	"log"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

type customHttpTransport struct {
	shopifyToken string
}

func (t *customHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(service.HeaderKeyUserAgent, service.HeaderValueUserAgent)
	req.Header.Add(service.HeaderKeyAccessToken, t.shopifyToken)
	return http.DefaultTransport.RoundTrip(req)
}

type Client struct {
	graphql graphql.Client
}

func NewClient() *Client {
	endpoint := fmt.Sprintf("https://%v/admin/api/%v/graphql.json", service.ShopId, service.ShopifyGraphqlApiVersion)
	httpClient := &http.Client{Transport: &customHttpTransport{shopifyToken: service.AccessToken}}
	return &Client{
		graphql: graphql.NewClient(endpoint, httpClient),
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
