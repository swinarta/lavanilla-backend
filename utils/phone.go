package utils

import "strings"

func NormalizePhone(phone string) string {
	switch {
	case strings.HasPrefix(phone, "0"):
		return "+62" + phone[1:]
	case strings.HasPrefix(phone, "62"):
		return "+" + phone
	}
	return phone
}
