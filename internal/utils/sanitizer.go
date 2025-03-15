package utils

import (
	"html"
	"regexp"
	"strings"
)

func SanitizeString(input string) string {
	safe := html.EscapeString(input)

	re := regexp.MustCompile(`(?i)<script.*?>.*?</script>`)
	safe = re.ReplaceAllString(safe, "")

	reIframe := regexp.MustCompile(`(?i)<iframe.*?>.*?</iframe>`)
	safe = reIframe.ReplaceAllString(safe, "")

	reObject := regexp.MustCompile(`(?i)<object.*?>.*?</object>`)
	safe = reObject.ReplaceAllString(safe, "")

	reEmbed := regexp.MustCompile(`(?i)<embed.*?>.*?</embed>`)
	safe = reEmbed.ReplaceAllString(safe, "")

	sqlPatterns := []string{"'", "--", ";--", "/*", "*/", "@@", "@"}
	for _, pattern := range sqlPatterns {
		safe = strings.ReplaceAll(safe, pattern, "")
	}

	return safe
}
