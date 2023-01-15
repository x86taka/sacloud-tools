package utils

import "strings"

func FormatHCL(name string) string {
	result := strings.ReplaceAll(name, " ", "_")
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, ".", "")
	result = strings.ToLower(result)
	return result
}
