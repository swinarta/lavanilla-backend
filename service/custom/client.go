package custom

import (
	"context"
	"fmt"
	"lavanilla/graphql/self-service/model"
	"lavanilla/service"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"github.com/samber/lo"
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
