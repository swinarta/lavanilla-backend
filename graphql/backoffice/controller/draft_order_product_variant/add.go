package draft_order_product_variant

import (
	"context"
	"errors"
	"lavanilla/graphql/backoffice/model"
	"lavanilla/service/custom"
	"lavanilla/service/shopify"

	"github.com/samber/lo"
)

func (h *Handler) Add(ctx context.Context, draftOrderID string, variantID string, quantity int) (*model.Order, error) {
	order, err := h.shopifyClient.GetDraftOrder(ctx, draftOrderID)
	if err != nil {
		return nil, err
	}
	_, found := lo.Find(order.DraftOrder.LineItems.Nodes, func(item shopify.GetDraftOrderDraftOrderLineItemsDraftOrderLineItemConnectionNodesDraftOrderLineItem) bool {
		return item.Variant.Id == variantID
	})

	if found {
		return nil, errors.New("variant already exists")
	}

	if err = h.shopifyClient.CheckDraftOrderStartedByDesigner(ctx, draftOrderID); err != nil {
		return nil, err
	}

	existingLineItems := lo.Map(order.DraftOrder.LineItems.Nodes, func(item shopify.GetDraftOrderDraftOrderLineItemsDraftOrderLineItemConnectionNodesDraftOrderLineItem, _ int) custom.DraftOrderLineItemInput {
		return custom.DraftOrderLineItemInput{
			Quantity:  item.Quantity,
			VariantId: item.Variant.Id,
		}
	})

	newLineItems := append(existingLineItems, custom.DraftOrderLineItemInput{
		Quantity:  quantity,
		VariantId: variantID,
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
