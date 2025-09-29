package draft_order_product_variant

import (
	"context"
	"errors"
	"lavanilla/graphql/backoffice/model"
	"lavanilla/service/custom"
	"lavanilla/service/shopify"

	"github.com/samber/lo"
)

func (h *Handler) Update(ctx context.Context, draftOrderID string, variantID string, quantity int) (*model.Order, error) {
	if err := h.shopifyClient.CheckDraftOrderStartedByDesigner(ctx, draftOrderID); err != nil {
		return nil, err
	}

	order, err := h.shopifyClient.GetDraftOrder(ctx, draftOrderID)
	if err != nil {
		return nil, err
	}
	_, found := lo.Find(order.DraftOrder.LineItems.Nodes, func(item shopify.GetDraftOrderDraftOrderLineItemsDraftOrderLineItemConnectionNodesDraftOrderLineItem) bool {
		return item.Variant.Id == variantID
	})

	if !found && quantity > 0 {
		return nil, errors.New("variant does not exist")
	}

	if !found && quantity == 0 {
		return &model.Order{
			ID:   order.DraftOrder.Id,
			Name: order.DraftOrder.Name,
		}, nil
	}

	newLineItems := lo.FilterMap(order.DraftOrder.LineItems.Nodes, func(item shopify.GetDraftOrderDraftOrderLineItemsDraftOrderLineItemConnectionNodesDraftOrderLineItem, _ int) (custom.DraftOrderLineItemInput, bool) {
		if item.Variant.Id == variantID && quantity == 0 {
			return custom.DraftOrderLineItemInput{
				Quantity:  quantity,
				VariantId: item.Variant.Id,
			}, false
		}

		if item.Variant.Id == variantID {
			return custom.DraftOrderLineItemInput{
				Quantity:  quantity,
				VariantId: item.Variant.Id,
			}, true
		}
		return custom.DraftOrderLineItemInput{
			Quantity:  item.Quantity,
			VariantId: item.Variant.Id,
		}, true
	})

	_, err = h.customClient.DraftOrderUpdateLineItems(ctx, draftOrderID, newLineItems)
	if err != nil {
		return nil, err
	}

	return &model.Order{
		ID:   order.DraftOrder.Id,
		Name: order.DraftOrder.Name,
	}, nil
}
