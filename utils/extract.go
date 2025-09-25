package utils

import (
	"fmt"
	"strings"
)

func ExtractIDWithPrefix(gid string, prefix string) (string, string, error) {
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
