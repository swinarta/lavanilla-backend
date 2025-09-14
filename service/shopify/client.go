package shopify

import (
	"bytes"
	"context"
	"fmt"
	"io"
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

func (c *Client) CreateDraftOrder1(ctx context.Context) (bool, error) {
	// c.graphql.MakeRequest(ctx)
	return true, nil
}

func (c *Client) CreateDraftOrder(ctx context.Context) (*CreateDraftOrderResponse, error) {
	return CreateDraftOrder(ctx, c.graphql, DraftOrderInput{
		Email: "sastraw@gmail.com",
		LineItems: []DraftOrderLineItemInput{
			{
				VariantId: "gid://shopify/ProductVariant/45539712303303",
				Quantity:  2,
				// OriginalUnitPriceWithCurrency: MoneyInput{
				// 	Amount:       "1000",
				// 	CurrencyCode: CurrencyCodeIdr,
				// },
				// AppliedDiscount: DraftOrderAppliedDiscountInput{
				// 	AmountWithCurrency: MoneyInput{
				// 		Amount:       "1000",
				// 		CurrencyCode: CurrencyCodeIdr,
				// 	},
				// 	ValueType: DraftOrderAppliedDiscountTypeFixedAmount,
				// },
				// PriceOverride: MoneyInput{
				// 	Amount:       "0",
				// 	CurrencyCode: CurrencyCodeIdr,
				// },
			},
		},
		// PresentmentCurrencyCode: CurrencyCodeIdr,
		// BillingAddress: MailingAddressInput{
		// 	CountryCode: CountryCodeId,
		// },
		// ShippingAddress: MailingAddressInput{
		// 	CountryCode: CountryCodeId,
		// },
		// ShippingLine: ShippingLineInput{
		// 	PriceWithCurrency: MoneyInput{
		// 		Amount:       "0",
		// 		CurrencyCode: CurrencyCodeIdr,
		// 	},
		// },
		// PaymentTerms: PaymentTermsInput{
		// 	PaymentTermsTemplateId: "gid://shopify/PaymentTermsTemplate/1",
		// },
		// PurchasingEntity: PurchasingEntityInput{
		// 	CustomerId: "gid://shopify/Customer/8185109446855",
		// 	// PurchasingCompany: PurchasingCompanyInput{
		// 	// 	CompanyContactId:  "gid://shopify/CompanyContact/10443701868879",
		// 	// 	CompanyId:         "gid://shopify/Company/10443701868879",
		// 	// 	CompanyLocationId: "gid://shopify/CompanyLocation/10443701868879",
		// 	// },
		// },
		// AppliedDiscount: DraftOrderAppliedDiscountInput{
		// 	AmountWithCurrency: MoneyInput{
		// 		Amount:       "0",
		// 		CurrencyCode: CurrencyCodeIdr,
		// 	},
		// 	Value:     0,
		// 	ValueType: DraftOrderAppliedDiscountTypeFixedAmount,
		// },
	})
}
