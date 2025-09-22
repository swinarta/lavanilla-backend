package shopify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lavanilla/graphql/backoffice/model"
	"lavanilla/service"
	"lavanilla/service/metadata"
	"log"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

const NameSpace = "LVN-APP"

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

func (c *Client) GetCustomer(ctx context.Context, email *string, phone *string) (*GetCustomerResponse, error) {
	query := ""
	if email != nil && phone != nil {
		query = fmt.Sprintf("email:%s OR phone:%s", *email, *phone)
	}
	if email != nil && phone == nil {
		query = fmt.Sprintf("email:%s", *email)
	}
	if phone != nil && email == nil {
		query = fmt.Sprintf("phone:%s", *phone)
	}
	return GetCustomer(ctx, c.graphql, query)
}

func (c *Client) CreateCustomer(ctx context.Context, email *string, phone *string) (*CreateCustomerResponse, error) {
	input := CustomerInput{}
	if email != nil {
		input.Email = *email
	}
	if phone != nil {
		input.Phone = *phone
	}
	return CreateCustomer(ctx, c.graphql, input)
}

func (c *Client) DraftOrderComplete(ctx context.Context, id string) (*DraftOrderCompleteResponse, error) {
	return DraftOrderComplete(ctx, c.graphql, id)
}

func (c *Client) DraftOrderUpdate(ctx context.Context, id string, input DraftOrderInput) (*DraftOrderUpdateResponse, error) {
	return DraftOrderUpdate(ctx, c.graphql, id, input)
}

func (c *Client) GetDraftOrders(ctx context.Context, tag *string, status *model.DraftOrderStatus) (*GetDraftOrdersResponse, error) {
	query := ""
	if tag != nil {
		query = fmt.Sprintf("tag:%s", *tag)
	}
	if status != nil {
		if query != "" {
			query += " AND "
		}
		if *status == model.DraftOrderStatusOpen {
			query += "status:OPEN"
		}
		if *status == model.DraftOrderStatusCompleted {
			query += "status:COMPLETED"
		}
	}
	return GetDraftOrders(ctx, c.graphql, query)
}

func (c *Client) GetDraftOrder(ctx context.Context, id string) (*GetDraftOrderResponse, error) {
	return GetDraftOrder(ctx, c.graphql, id)
}

func (c *Client) MetaDataAdd(ctx context.Context, ownerId string, key metadata.KeyName, value []byte) (*MetaDataAddResponse, error) {
	return MetaDataAdd(ctx, c.graphql, ownerId, NameSpace, key, string(value))
}

func (c *Client) GetDraftOrderMetaField(ctx context.Context, orderId string, key string) (*GetDraftOrderMetaFieldResponse, error) {
	return GetDraftOrderMetaField(ctx, c.graphql, orderId, NameSpace, key)
}

func (c *Client) TimestampAdd(ctx context.Context, orderId string, value metadata.Timeline) (*string, error) {
	found, err := GetDraftOrderMetaField(ctx, c.graphql, orderId, NameSpace, metadata.TImeLineKeyName)
	if err != nil {
		return nil, err
	}
	var payload []metadata.Timeline
	if found == nil {
		payload = []metadata.Timeline{value}
	}
	marshal, _ := json.Marshal(payload)
	add, err := MetaDataAdd(ctx, c.graphql, orderId, NameSpace, metadata.TImeLineKeyName, string(marshal))
	if err != nil {
		return nil, err
	}
	if len(add.MetafieldsSet.UserErrors) > 0 {
		return nil, errors.New(string(add.MetafieldsSet.UserErrors[0].Code))
	}
	return &add.MetafieldsSet.Metafields[0].Id, nil
}
