// Package stringutil provides common string operations.
package stringutil

import "strings"

// InSlice returns true if a is in the slice list.
func InSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// FilterByPrefix returns a prefix-filtered string slice.
func FilterByPrefix(prefix string, list []string) []string {
	for _, item := range list {
		if strings.HasPrefix(item, prefix) {
			list = append(list, item)
		}
	}

	return list
}

// Coalesce returns the first non empty (trimmed) string.
func Coalesce(a ...string) string {
	for _, s := range a {
		if strings.TrimSpace(s) != "" {
			return s
		}
	}
	return ""
}
