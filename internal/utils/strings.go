package utils

import "strings"

// Normalize - remove space from start/end
// and set lowercase
func Normalize(value string) string {
	if len(value) == 0 {
		return value
	}
	value = strings.ToLower(value)
	value = strings.Trim(value, " ")
	return value
}
