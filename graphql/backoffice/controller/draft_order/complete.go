package draft_order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lavanilla/graphql/backoffice/model"
	"lavanilla/service"
	"lavanilla/service/metadata"
	S3util "lavanilla/service/s3util"
	"lavanilla/utils"
	"time"

	"github.com/samber/lo"
)

func (h *Handler) Complete(ctx context.Context, draftOrderID string) (*model.Order, error) {

	if err := h.shopifyClient.CheckDraftOrderStartedByDesigner(ctx, draftOrderID); err != nil {
		return nil, err
	}

	tag, err := h.shopifyClient.RemoveTag(ctx, draftOrderID, metadata.DesignerInProgressKeyName)
	if err != nil {
		return nil, err
	}
	if len(tag.TagsRemove.UserErrors) > 0 {
		return nil, errors.New(tag.TagsRemove.UserErrors[0].Message)
	}

	order, err := h.shopifyClient.DraftOrderComplete(ctx, draftOrderID)
	if err != nil {
		return nil, err
	}
	if len(order.DraftOrderComplete.UserErrors) > 0 {
		return nil, errors.New(order.DraftOrderComplete.UserErrors[0].Message)
	}
	newOrderId, newGlobalOrderId, _ := utils.ExtractID(order.DraftOrderComplete.DraftOrder.Order.Id)

	// rename from draftOrder.name to order.id
	if err := S3util.RenameS3Directory(ctx, h.s3Client, service.S3BucketOrder, fmt.Sprintf("%s/", order.DraftOrderComplete.DraftOrder.Name), fmt.Sprintf("%s/", newOrderId)); err != nil {
		return nil, err
	}

	// add draft order timeline and copy timeline to order
	now := time.Now()
	designerEndMetadata := metadata.Timeline{
		Timestamp: now,
		Action:    "DESIGNER_END",
	}
	_, newMetadata, err := h.shopifyClient.TimestampAdd(ctx, draftOrderID, designerEndMetadata)
	if err != nil {
		return nil, err
	}
	if _, err := h.shopifyClient.OrderTimestampInit(ctx, newGlobalOrderId, newMetadata); err != nil {
		return nil, err
	}

	// set designer perf metadata to draft order
	var designerJob *metadata.DesignerJob
	if _, err = h.shopifyClient.GetDraftOrderMetaField(ctx, draftOrderID, metadata.DesignerKeyName, &designerJob); err != nil {
		return nil, err
	}

	if designerJob != nil {
		designerJob.EndAt = lo.ToPtr(now)
	}

	// TODO: goroutine
	marshal, _ := json.Marshal(designerJob)
	if _, err = h.shopifyClient.MetaDataAdd(ctx, draftOrderID, metadata.DesignerKeyName, marshal); err != nil {
		return nil, err
	}
	if _, err = h.shopifyClient.MetaDataAdd(ctx, newGlobalOrderId, metadata.DesignerKeyName, marshal); err != nil {
		return nil, err
	}

	// TODO: use get?
	return &model.Order{
		ID:        newGlobalOrderId,
		Name:      order.DraftOrderComplete.DraftOrder.Order.Name,
		LineItems: nil,
		Timelines: nil,
		CreatedAt: order.DraftOrderComplete.DraftOrder.CreatedAt,
		Customer: &model.Customer{
			Name: order.DraftOrderComplete.DraftOrder.Customer.DisplayName,
		},
	}, nil
}
