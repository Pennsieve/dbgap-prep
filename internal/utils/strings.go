package utils

import "strings"

// EqualIgnoreCaseTrimmed returns true if a and b are equal ignoring case after having been trimmed of leading and trailing whitespace.
// It returns false otherwise.
func EqualIgnoreCaseTrimmed(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}
