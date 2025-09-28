package draft_order

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"lavanilla/service"
	"lavanilla/utils"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (h *Handler) DownloadAssetsDesigner(ctx context.Context, draftOrderID string) (*string, error) {
	draftOrderID, globalDraftOrderId, err := utils.ExtractIDWithDraftOrderPrefix(draftOrderID)
	if err != nil {
		return nil, err
	}

	if err = h.shopifyClient.CheckDraftOrderStartedByDesigner(ctx, globalDraftOrderId); err != nil {
		return nil, err
	}

	draftOrder, err := h.shopifyClient.GetDraftOrder(ctx, globalDraftOrderId)
	if err != nil {
		return nil, err
	}
	orderName := draftOrder.DraftOrder.Name

	resp, err := h.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(service.S3BucketSelfService),
		Prefix: aws.String(draftOrderID),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to list objects: %v", err))
	}

	if len(resp.Contents) <= 0 {
		return nil, errors.New("no assets to download")
	}

	results := make(chan utils.FileData)
	var wg sync.WaitGroup
	for _, content := range resp.Contents {
		wg.Add(1)

		go func(content types.Object) {
			defer wg.Done()
			obj, err := h.s3Client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: aws.String(service.S3BucketSelfService),
				Key:    content.Key,
			})
			if err != nil {
				log.Printf("failed to get object %s: %v\n", *content.Key, err)
				results <- utils.FileData{Key: *content.Key, Err: err}
				return
			}
			defer obj.Body.Close()

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, obj.Body)
			if err != nil {
				results <- utils.FileData{Key: *content.Key, Err: err}
				return
			}
			filename := *content.Key
			parts := strings.Split(filename, "/")
			if len(parts) > 0 {
				parts[0] = orderName
			}
			results <- utils.FileData{Key: strings.Join(parts, "/"), Data: buf.Bytes()}
		}(content)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	urlResp, err := utils.CreateZipArchive(ctx, orderName, h.s3Client, h.s3PresignClient, results)
	if err != nil {
		return nil, fmt.Errorf("failed to create zip archive: %w", err)
	}

	return urlResp, nil
}
