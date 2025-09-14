package custom

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"lavanilla/graphql/self-service/model"
	"log"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"github.com/samber/lo"
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
	if req.Body != nil {
		// Read and copy the body
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("error reading request body: %v", err)
		} else {
			log.Printf("Request Body: %s", string(bodyBytes))
		}

		// Reconstruct the Body so the transport can still use it
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
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

func (c *Client) CreateDraftOrder(ctx context.Context, input model.OrderInput) (*CreateDraftOrderResponse, error) {
	return CreateDraftOrder(ctx, c.graphql, DraftOrderInput{
		Email: "sastraw@gmail.com",
		LineItems: lo.Map(input.Items, func(item *model.LineItem, _ int) DraftOrderLineItemInput {
			return DraftOrderLineItemInput{
				VariantId: item.VariantID,
				Quantity:  item.Quantity,
			}
		}),
	})
}
