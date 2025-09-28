package draft_order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lavanilla/service"
	"lavanilla/service/metadata"
	S3util "lavanilla/service/s3util"
	"lavanilla/utils"
	"time"

	"github.com/samber/lo"
)

func (h *Handler) Complete(ctx context.Context, id string) (bool, error) {

	if err := h.shopifyClient.CheckDraftOrderStartedByDesigner(ctx, id); err != nil {
		return false, err
	}

	tag, err := h.shopifyClient.RemoveTag(ctx, id, metadata.DesignerInProgressKeyName)
	if err != nil {
		return false, err
	}
	if len(tag.TagsRemove.UserErrors) > 0 {
		return false, errors.New(tag.TagsRemove.UserErrors[0].Message)
	}

	order, err := h.shopifyClient.DraftOrderComplete(ctx, id)
	if err != nil {
		return false, err
	}

	newOrderId, _, err := utils.ExtractID(order.DraftOrderComplete.DraftOrder.Order.Id)
	if err != nil {
		return false, err
	}

	// rename from draftOrder.name to order.id
	if err := S3util.RenameS3Directory(ctx, h.s3Client, service.S3BucketOrder, fmt.Sprintf("%s/", order.DraftOrderComplete.DraftOrder.Name), fmt.Sprintf("%s/", newOrderId)); err != nil {
		return false, err
	}

	now := time.Now()
	_, err = h.shopifyClient.TimestampAdd(ctx, id, metadata.Timeline{
		Timestamp: now,
		Action:    "DESIGNER_END",
	})
	if err != nil {
		return false, err
	}
	var designerJob *metadata.DesignerJob
	_, err = h.shopifyClient.GetDraftOrderMetaField(ctx, id, metadata.DesignerKeyName, &designerJob)
	if err != nil {
		return false, err
	}

	if designerJob != nil {
		designerJob.EndAt = lo.ToPtr(now)
	}

	marshal, _ := json.Marshal(designerJob)
	_, err = h.shopifyClient.MetaDataAdd(ctx, id, metadata.DesignerKeyName, marshal)
	if err != nil {
		return false, err
	}

	return true, nil
}
