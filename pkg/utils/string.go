package utils

import "unicode/utf8"

// LimitString makes sure that a is not bigger than limit or as long as limit with the appended Unicode that looks like three dots
func LimitString(a string, limit int) string {
	if utf8.RuneCountInString(a) > limit {
		return a[0:limit-1] + "…"
	}
	return a
}
