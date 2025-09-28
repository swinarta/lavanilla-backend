package custom

import (
	"context"
	"encoding/json"
	"fmt"
	"lavanilla/graphql/self-service/model"
	"lavanilla/service"
	"lavanilla/service/metadata"
	"net/http"
	"time"

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
	create := []metadata.Timeline{{
		Timestamp: time.Now(),
		Action:    "DRAFT_ORDER_CREATED",
	}}
	marshal, _ := json.Marshal(create)
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
			"LVN-APP",
		},
		PaymentTerms: PaymentTermsInput{
			PaymentTermsTemplateId: "gid://shopify/PaymentTermsTemplate/1",
		},
		Metafields: []MetafieldInput{{
			// Id:        "",
			Namespace: "LVN-APP",
			Key:       "TIMELINE",
			Value:     string(marshal),
			Type:      "json",
		}},
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

func (c *Client) DraftOrderUpdateLineItems(ctx context.Context, id string, input []DraftOrderLineItemInput) (*DraftOrderUpdateLineItemsResponse, error) {
	return DraftOrderUpdateLineItems(ctx, c.graphql, id, input)
}
