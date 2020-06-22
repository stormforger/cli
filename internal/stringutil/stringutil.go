// Package stringutil provides common string operations.
package stringutil

// InSlice returns true if a is in the slice list.
func InSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
