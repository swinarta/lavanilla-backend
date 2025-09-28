package order

import (
	"context"
	"lavanilla/graphql/backoffice/model"
	"lavanilla/service/metadata"
	"lavanilla/service/shopify"

	"github.com/99designs/gqlgen/graphql"
	"github.com/samber/lo"
)

func (h *Handler) Order(ctx context.Context, orderID string) (*model.Order, error) {
	order, err := h.shopifyClient.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	fields := graphql.CollectFieldsCtx(ctx, nil)
	collectTimelines := false
	for _, field := range fields {
		if field.Name == "timelines" {
			collectTimelines = true
		}
	}

	var timelineData []metadata.Timeline
	if collectTimelines {
		// timelineData, err = h.shopifyClient.GetDraftOrderTimeline(ctx, draftOrderID)
		// if err != nil {
		// 	return nil, err
		// }
	}

	return &model.Order{
		ID:        order.Order.Id,
		Name:      order.Order.Name,
		CreatedAt: order.Order.CreatedAt,
		Customer: &model.Customer{
			Name: order.Order.Customer.DisplayName,
		},
		Timelines: lo.Map(timelineData, func(item metadata.Timeline, _ int) *model.Timeline {
			return &model.Timeline{
				EventAt: item.Timestamp,
				Action:  model.EventAction(item.Action),
			}
		}),

		LineItems: lo.Map(order.Order.LineItems.Nodes, func(item shopify.GetOrderOrderLineItemsLineItemConnectionNodesLineItem, _ int) *model.LineItem {
			return &model.LineItem{
				Quantity: item.Quantity,
				Product: &model.Product{
					ID:    item.Product.Id,
					Title: item.Product.Title,
				},
				Variant: &model.ProductVariant{
					ID:  item.Id,
					Sku: item.Sku,
				},
			}
		}),
	}, nil
}
