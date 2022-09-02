package stringsx

import "strings"

func IsEmpty(s string) bool {
	var empty string
	return s == empty
}

func IsEmptyWithTrim(s string) bool {
	var empty string
	return strings.TrimSpace(s) == empty
}
