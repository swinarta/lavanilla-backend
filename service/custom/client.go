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

func (c *Client) CreateDraftOrder(ctx context.Context, input model.OrderInput, customerId string) (*CreateDraftOrderResponse, error) {
	payload := DraftOrderInput{
		PurchasingEntity: PurchasingEntityInput{CustomerId: customerId},
		LineItems: lo.Map(input.Items, func(item *model.LineItem, _ int) DraftOrderLineItemInput {
			return DraftOrderLineItemInput{
				VariantId: item.VariantID,
				Quantity:  item.Quantity,
			}
		}),
		Tags: []string{
			"DESAINER",
		},
		PaymentTerms: PaymentTermsInput{
			PaymentTermsTemplateId: "gid://shopify/PaymentTermsTemplate/1",
		},
	}
	if input.Customer.Email != nil {
		payload.Email = *input.Customer.Email
	}
	if input.Customer.Phone != nil {
		payload.Phone = *input.Customer.Phone
	}
	if input.Note != nil {
		payload.Note = *input.Note
	}
	return CreateDraftOrder(ctx, c.graphql, payload)
}
