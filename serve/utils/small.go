package utils

import (
	mapset "github.com/deckarep/golang-set/v2"
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

func HasEmpty(vs ...string) bool {
	if vs == nil {
		return true
	}
	for _, v := range vs {
		if len(v) == 0 {
			return true
		}
	}
	return false
}

func ConvertSetToArray[T comparable](set mapset.Set[T]) []T {
	var array []T
	set.Each(func(e T) bool {
		array = append(array, e)
		return false
	})
	return array
}
