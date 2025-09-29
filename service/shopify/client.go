package shopify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lavanilla/graphql/backoffice/model"
	"lavanilla/service"
	"lavanilla/service/metadata"
	"net/http"
	"reflect"

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

func (c *Client) AddTag(ctx context.Context, orderId string, tag metadata.KeyTag) (*TagsAddResponse, error) {
	return TagsAdd(ctx, c.graphql, orderId, tag)
}

func (c *Client) RemoveTag(ctx context.Context, orderId string, tag metadata.KeyTag) (*TagsRemoveResponse, error) {
	return TagsRemove(ctx, c.graphql, orderId, tag)
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

func (c *Client) CreateCustomer(ctx context.Context, name string, email *string, phone *string) (*CreateCustomerResponse, error) {
	input := CustomerInput{
		FirstName: name,
	}
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

func (c *Client) GetOrders(ctx context.Context) (*GetOrdersResponse, error) {
	query := "tag:LVN-APP AND financial_status:paid"
	// if tag != nil {
	// 	query = fmt.Sprintf("tag:%s", *tag)
	// }
	// if status != nil {
	// 	if query != "" {
	// 		query += " AND "
	// 	}
	// 	if *status == model.DraftOrderStatusOpen {
	// 		query += "status:OPEN"
	// 	}
	// 	if *status == model.DraftOrderStatusCompleted {
	// 		query += "status:COMPLETED"
	// 	}
	// }
	return GetOrders(ctx, c.graphql, query)
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

func (c *Client) GetOrder(ctx context.Context, id string) (*GetOrderResponse, error) {
	return GetOrder(ctx, c.graphql, id)
}

func (c *Client) GetDraftOrder(ctx context.Context, id string) (*GetDraftOrderResponse, error) {
	return GetDraftOrder(ctx, c.graphql, id)
}

func (c *Client) MetaDataAdd(ctx context.Context, ownerId string, key metadata.KeyName, value []byte) (*MetaDataAddResponse, error) {
	return MetaDataAdd(ctx, c.graphql, ownerId, NameSpace, key, string(value))
}

func (c *Client) DraftOrderCustomAttributesAdd(ctx context.Context, draftOrderId string, key metadata.AttrKey, value string) (*DraftOrderCustomAttributesAddResponse, error) {
	return DraftOrderCustomAttributesAdd(ctx, c.graphql, draftOrderId, key, value)
}

func (c *Client) GetDraftOrderMetaField(ctx context.Context, orderId string, key string, value any) (*GetDraftOrderMetaFieldResponse, error) {
	res, err := GetDraftOrderMetaField(ctx, c.graphql, orderId, NameSpace, key)
	if err != nil {
		return nil, err
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		fmt.Println("value must be a non-nil pointer")
		return res, nil
	}

	elem := rv.Elem()
	newVal := reflect.New(elem.Type().Elem())
	_ = json.Unmarshal([]byte(res.DraftOrder.Metafield.Value), newVal.Interface())
	elem.Set(newVal)

	return res, nil
}

func (c *Client) CheckDraftOrderStartedByDesigner(ctx context.Context, draftOrderID string) error {
	var designerJob *metadata.DesignerJob
	if _, err := c.GetDraftOrderMetaField(ctx, draftOrderID, metadata.DesignerKeyName, &designerJob); err != nil {
		return err
	}

	if designerJob == nil || designerJob.StartAt == nil {
		return errors.New("job not started")
	}

	return nil
}

func (c *Client) GetDraftOrderTimeline(ctx context.Context, orderId string) ([]metadata.Timeline, error) {
	found, err := GetDraftOrderMetaField(ctx, c.graphql, orderId, NameSpace, metadata.TImeLineKeyName)
	if err != nil {
		return nil, err
	}
	var payload []metadata.Timeline
	if found.DraftOrder.Metafield.Value != "" {
		if err := json.Unmarshal([]byte(found.DraftOrder.Metafield.Value), &payload); err != nil {
			return nil, err
		}
	}
	return payload, nil
}

func (c *Client) TimestampAdd(ctx context.Context, orderId string, value metadata.Timeline) (*string, []metadata.Timeline, error) {
	payload, err := c.GetDraftOrderTimeline(ctx, orderId)
	payload = append(payload, value)

	marshal, _ := json.Marshal(payload)
	add, err := MetaDataAdd(ctx, c.graphql, orderId, NameSpace, metadata.TImeLineKeyName, string(marshal))
	if err != nil {
		return nil, nil, err
	}
	if len(add.MetafieldsSet.UserErrors) > 0 {
		return nil, nil, errors.New(string(add.MetafieldsSet.UserErrors[0].Code))
	}
	return &add.MetafieldsSet.Metafields[0].Id, payload, nil
}

func (c *Client) GetOrderMetaField(ctx context.Context, orderId string, key string, value any) (*GetOrderMetaFieldResponse, error) {
	res, err := GetOrderMetaField(ctx, c.graphql, orderId, NameSpace, key)
	if err != nil {
		return nil, err
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		fmt.Println("value must be a non-nil pointer")
		return res, nil
	}

	elem := rv.Elem()
	newVal := reflect.New(elem.Type().Elem())
	_ = json.Unmarshal([]byte(res.Order.Metafield.Value), newVal.Interface())
	elem.Set(newVal)

	return res, nil
}

func (c *Client) GetOrderTimeline(ctx context.Context, orderId string) ([]metadata.Timeline, error) {
	found, err := GetOrderMetaField(ctx, c.graphql, orderId, NameSpace, metadata.TImeLineKeyName)
	if err != nil {
		return nil, err
	}
	var payload []metadata.Timeline
	if found.Order.Metafield.Value != "" {
		if err := json.Unmarshal([]byte(found.Order.Metafield.Value), &payload); err != nil {
			return nil, err
		}
	}
	return payload, nil
}

func (c *Client) OrderTimestampInit(ctx context.Context, orderId string, payload []metadata.Timeline) (*string, error) {
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

func (c *Client) OrderTimestampAdd(ctx context.Context, orderId string, value metadata.Timeline) (*string, error) {
	payload, err := c.GetDraftOrderTimeline(ctx, orderId)
	payload = append(payload, value)

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
