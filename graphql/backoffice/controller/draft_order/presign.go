package draft_order

import (
	"context"
	"fmt"
	"lavanilla/service"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (h *Handler) PresignedURLDesigner(ctx context.Context, orderName string, sku string, qty int) ([]string, error) {
	var result []string
	for i := 0; i < qty; i++ {
		filename := fmt.Sprintf("%s/%s/%d.jpeg", orderName, sku, time.Now().Unix()+1)
		object, err := h.s3PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(service.S3BucketOrder),
			ContentType: aws.String("image/jpeg"),
			Key:         aws.String(filename),
		}, func(options *s3.PresignOptions) {
			options.Expires = 15 * time.Minute
		})
		if err != nil {
			return nil, err
		}
		result = append(result, object.URL)
	}
	return result, nil
}
