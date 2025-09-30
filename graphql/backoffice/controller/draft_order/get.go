package draft_order

import (
	"context"
	"errors"
	"fmt"
	"lavanilla/graphql/backoffice/model"
	"lavanilla/service"
	"lavanilla/service/metadata"
	"lavanilla/service/shopify"
	"lavanilla/utils"
	"net/url"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/samber/lo"
)

func (h *Handler) DraftOrder(ctx context.Context, draftOrderID string) (*model.Order, error) {
	_, draftOrderID, err := utils.ExtractIDWithDraftOrderPrefix(draftOrderID)
	if err != nil {
		return nil, err
	}

	order, err := h.shopifyClient.GetDraftOrder(ctx, draftOrderID)
	if err != nil {
		return nil, err
	}

	s3resp, err := h.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(service.S3BucketOrder),
		Prefix: aws.String(order.DraftOrder.Name),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to list objects: %v", err))
	}

	objectMap := make(map[string][]string)
	for _, content := range s3resp.Contents {
		parts := strings.Split(*content.Key, "/")
		if len(parts) < 2 {
			continue
		}
		key := parts[1]
		objectMap[key] = append(objectMap[key], *content.Key)
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
		timelineData, err = h.shopifyClient.GetDraftOrderTimeline(ctx, draftOrderID)
		if err != nil {
			return nil, err
		}
	}

	return &model.Order{
		ID:   order.DraftOrder.Id,
		Name: order.DraftOrder.Name,
		LineItems: lo.Map(order.DraftOrder.LineItems.Nodes, func(item shopify.GetDraftOrderDraftOrderLineItemsDraftOrderLineItemConnectionNodesDraftOrderLineItem, _ int) *model.LineItem {
			designerNote, foundDesignerNote := lo.Find(item.CustomAttributes, func(item shopify.GetDraftOrderDraftOrderLineItemsDraftOrderLineItemConnectionNodesDraftOrderLineItemCustomAttributesAttribute) bool {
				return item.Key == metadata.DesignerNoteKeyName
			})

			foundImages, _ := objectMap[item.Sku]
			variantPrice, _ := strconv.ParseFloat(item.Variant.Price, 64)
			return &model.LineItem{
				Product: &model.Product{
					ID:    item.Id,
					Title: item.Title,
				},
				Quantity: item.Quantity,
				Variant: &model.ProductVariant{
					ID:    item.Variant.Id,
					Title: item.Variant.Title,
					Sku:   item.Variant.Sku,
					Price: variantPrice,
					Image: []string{item.Variant.Image.Url},
				},
				UploadedImages: lo.Map(foundImages, func(item string, _ int) string {
					return fmt.Sprintf("%s/%s", service.CdnOrder, url.QueryEscape(item))
				}),
				DesignerNote: lo.If(foundDesignerNote, &designerNote.Value).Else(nil),
			}
		}),
		Timelines: lo.Map(timelineData, func(item metadata.Timeline, _ int) *model.Timeline {
			return &model.Timeline{
				EventAt: item.Timestamp,
				Action:  model.EventAction(item.Action),
			}
		}),
	}, nil
}
