package internal

import (
	"strings"
	"unicode"
)

func CleanStringData(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, s)
}

func TransformString(length int, s string) string {
	if s == "" {
		return "NULL"
	}
	transformed := s
	if strings.Contains(s, "'") {
		transformed = strings.ReplaceAll(s, "'", "''")
	}
	if len(transformed) > length {
		return "'" + transformed[:length] + "'"
	}
	return "'" + transformed + "'"
}
