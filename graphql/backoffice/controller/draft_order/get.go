package draft_order

import (
	"context"
	"lavanilla/graphql/backoffice/model"
	"lavanilla/service/shopify"

	"github.com/samber/lo"
)

func (h *Handler) Get(ctx context.Context, status *model.DraftOrderStatus) ([]*model.Order, error) {
	orders, err := h.shopifyClient.GetDraftOrders(ctx, lo.ToPtr("DESAINER"), status)
	if err != nil {
		return nil, err
	}
	return lo.Map(orders.DraftOrders.Nodes, func(item shopify.GetDraftOrdersDraftOrdersDraftOrderConnectionNodesDraftOrder, _ int) *model.Order {
		return &model.Order{
			ID:        item.Id,
			Name:      item.Name,
			CreatedAt: item.CreatedAt,
			Customer:  &model.Customer{Name: item.Customer.DisplayName},
		}
	}), nil
}
