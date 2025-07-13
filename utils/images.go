package utils

import (
	"encoding/base64"
	"strings"
)

func IsValidBase64Image(data string) bool {
	if !strings.HasPrefix(data, "data:image/") {
		return false
	}
	parts := strings.SplitN(data, ",", 2)
	if len(parts) != 2 {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(parts[1])
	return err == nil
}
