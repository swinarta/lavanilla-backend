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

func (c *Client) DraftOrderCreate(ctx context.Context, input model.OrderInput, customerId string) (*DraftOrderCreateResponse, error) {
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
	return DraftOrderCreate(ctx, c.graphql, payload)
}

func (c *Client) DraftOrderUpdate(ctx context.Context, id string, input DraftOrderInput) (*DraftOrderUpdateResponse, error) {
	return DraftOrderUpdate(ctx, c.graphql, id, input)
}

func (c *Client) DraftOrderUpdate1(ctx context.Context, id string) (*DraftOrderUpdateResponse, error) {
	input := DraftOrderInput{
		LineItems: []DraftOrderLineItemInput{
			{
				VariantId: "gid://shopify/DraftOrder/1001317793991",
				Quantity:  23,
			},
		},
	}
	return DraftOrderUpdate(ctx, c.graphql, id, input)
	// return DraftOrderUpdate1(ctx, c.graphql, id)
}

func (c *Client) DraftOrderUpdate2(ctx context.Context, id string, input []DraftOrderLineItemInput) (*DraftOrderUpdate2Response, error) {
	return DraftOrderUpdate2(ctx, c.graphql, id, input)
}
