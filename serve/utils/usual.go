package utils

import (
	"strings"
)

// Cut slices and returns the text before the prefix,
// between the prefix and suffix and back the suffix.
func Cut(s string, prefix string, suffix string) (before, middle, back string) {
	b1, a1, found := strings.Cut(s, suffix)
	if !found {
		return "", "", ""
	}
	b2, a2, found := strings.Cut(b1, prefix)
	if !found {
		return "", "", ""
	}
	return b2, a2, a1
}

func MapToArray[K int | int64 | string, V any](m map[K]V) []V {
	var arr []V
	for _, v := range m {
		arr = append(arr, v)
	}
	return arr
}
