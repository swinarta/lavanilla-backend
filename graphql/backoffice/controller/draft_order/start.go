package draft_order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lavanilla/service/metadata"
	"time"

	"github.com/samber/lo"
)

func (h *Handler) Start(ctx context.Context, id string) (bool, error) {
	var designerJob *metadata.DesignerJob
	if _, err := h.shopifyClient.GetDraftOrderMetaField(ctx, id, metadata.DesignerKeyName, &designerJob); err != nil {
		return false, err
	}

	if designerJob != nil && designerJob.StartAt != nil {
		return false, errors.New(fmt.Sprintf("meta fields %s are not empty", metadata.DesignerKeyName))
	}

	now := time.Now()
	job := metadata.DesignerJob{
		StartAt: lo.ToPtr(now),
		EndAt:   nil,
	}

	marshal, _ := json.Marshal(job)
	m, err := h.shopifyClient.MetaDataAdd(ctx, id, metadata.DesignerKeyName, marshal)
	if err != nil {
		return false, err
	}
	if len(m.MetafieldsSet.UserErrors) > 0 {
		return false, errors.New(string(m.MetafieldsSet.UserErrors[0].Code))
	}

	if _, _, err = h.shopifyClient.TimestampAdd(ctx, id, metadata.Timeline{
		Timestamp: now,
		Action:    "DESIGNER_START",
	}); err != nil {
		return false, err
	}

	tag, err := h.shopifyClient.AddTag(ctx, id, metadata.DesignerInProgressKeyName)
	if err != nil {
		return false, err
	}
	if len(tag.TagsAdd.UserErrors) > 0 {
		return false, errors.New(tag.TagsAdd.UserErrors[0].Message)
	}

	return true, nil
}
