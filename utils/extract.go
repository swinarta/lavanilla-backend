package utils

import (
	"errors"
	"fmt"
	"strings"
)

func GetGlobalDraftOrderId(draftOrderID string) string {
	return fmt.Sprintf("gid://shopify/DraftOrder/%s", draftOrderID)
}

func GetGlobalOrderId(orderID string) string {
	return fmt.Sprintf("gid://shopify/Order/%s", orderID)
}

func ExtractIDWithDraftOrderPrefix(gid string) (string, string, error) {
	return extractIDWithPrefix(gid, "gid://shopify/DraftOrder/")
}

func extractIDWithPrefix(gid string, prefix string) (string, string, error) {
	if !strings.HasPrefix(gid, prefix) {
		return "", "", errors.New("wrong prefix given")
	}
	parts := strings.Split(gid, "/")
	if len(parts) == 0 {
		return "", gid, fmt.Errorf("invalid gid: %s", gid)
	}
	return parts[len(parts)-1], gid, nil
}

func ExtractID(gid string) (string, string, error) {
	parts := strings.Split(gid, "/")
	if len(parts) == 0 {
		return "", gid, fmt.Errorf("invalid gid: %s", gid)
	}
	return parts[len(parts)-1], gid, nil
}
