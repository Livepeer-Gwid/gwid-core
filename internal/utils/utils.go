package utils

import (
	"strings"
	"unicode"
)

func ToKebabCase(str string) string {
	var result strings.Builder

	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 && result.Len() > 0 {
				str := result.String()
				if !strings.HasSuffix(str, "-") {
					result.WriteRune('-')
				}
			}
			result.WriteRune(unicode.ToLower(r))
		} else if unicode.IsSpace(r) || r == '_' {
			if result.Len() > 0 && !strings.HasSuffix(result.String(), "-") {
				result.WriteRune('-')
			}
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func ToSnakeCase(str string) string {
	var result strings.Builder
	runes := []rune(str)

	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				prevRune := runes[i-1]

				if unicode.IsLower(prevRune) || unicode.IsDigit(prevRune) {
					result.WriteRune('_')
				} else if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					result.WriteRune('_')
				}
			}

			result.WriteRune(unicode.ToLower(r))
		} else if unicode.IsSpace(r) || r == '-' {
			if result.Len() > 0 && !strings.HasSuffix(result.String(), "_") {
				result.WriteRune('_')
			}
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}

	return result.String()
}
