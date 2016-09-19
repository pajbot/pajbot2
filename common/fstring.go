package common

import "strings"

// FormatString formats a string with a map data
func FormatString(format string, data map[string]string) string {
	for k, v := range data {
		format = strings.Replace(format, "{"+k+"}", v, -1)
	}

	return format
}
