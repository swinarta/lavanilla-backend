package utils

import (
	"fmt"
	"strings"
)

func ExtractID(gid string) (string, error) {
	parts := strings.Split(gid, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid gid: %s", gid)
	}
	return parts[len(parts)-1], nil
}
