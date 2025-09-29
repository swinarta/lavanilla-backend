package draft_order_product_variant

import (
	"context"
	"errors"
	"lavanilla/service/metadata"
	"lavanilla/service/shopify"

	"github.com/samber/lo"
)

func (h *Handler) AddNote(ctx context.Context, draftOrderID string, variantID string, input string) (bool, error) {
	if err := h.shopifyClient.CheckDraftOrderStartedByDesigner(ctx, draftOrderID); err != nil {
		return false, err
	}

	order, err := h.shopifyClient.GetDraftOrder(ctx, draftOrderID)
	if err != nil {
		return false, err
	}
	foundVariant, found := lo.Find(order.DraftOrder.LineItems.Nodes, func(item shopify.GetDraftOrderDraftOrderLineItemsDraftOrderLineItemConnectionNodesDraftOrderLineItem) bool {
		return item.Variant.Id == variantID
	})

	if !found {
		return false, errors.New("variant does not exist")
	}

	if _, err := h.shopifyClient.DraftOrderLineItemCustomAttributesAdd(ctx, draftOrderID, variantID, foundVariant.Quantity, metadata.DesignerNoteKeyName, input); err != nil {
		return false, err
	}

	return true, nil
}
