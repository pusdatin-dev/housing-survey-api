package utils

import "strings"

func SplitAndTrim(types string, s string) []string {
	var result []string
	typeList := strings.Split(types, s)
	for _, item := range typeList {
		trimmedItem := strings.TrimSpace(item)
		if trimmedItem != "" {
			result = append(result, trimmedItem)
		}
	}
	return result
}
