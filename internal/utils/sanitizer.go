package utils

import (
	"html"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func SanitizeString(input string) string {
	safe := html.EscapeString(input)

	safe = regexp.MustCompile(`(?i)<script.*?>.*?</script>`).ReplaceAllString(safe, "")
	safe = regexp.MustCompile(`(?i)<iframe.*?>.*?</iframe>`).ReplaceAllString(safe, "")
	safe = regexp.MustCompile(`(?i)<object.*?>.*?</object>`).ReplaceAllString(safe, "")
	safe = regexp.MustCompile(`(?i)<embed.*?>.*?</embed>`).ReplaceAllString(safe, "")

	safe = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, safe)

	sqlPatterns := []string{"'", "--", ";--", "/*", "*/", "@@", "@", "DROP", "ALTER", "INSERT", "UPDATE", "DELETE", "SELECT"}
	for _, pattern := range sqlPatterns {
		safe = strings.ReplaceAll(safe, pattern, "")
	}

	return safe
}

func SanitizeNumberInt(value int) int {
	if value < 0 {
		return 0
	}
	return value
}

func SanitizeNumberFloat64(value float64) float64 {
	if math.IsNaN(value) {
		return 0
	}
	if value < 0 {
		return 0
	}
	return value
}

func SanitizeDate(date interface{}) time.Time {
	switch v := date.(type) {
	case string:
		parsedDate, err := time.Parse("2006-01-02", v)
		if err != nil {
			return time.Time{}
		}
		return parsedDate
	case time.Time:
		return v
	default:
		return time.Time{}
	}
}
