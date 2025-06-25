package helper

import (
	"regexp"
	"strings"
)

func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace underscores and spaces with dashes
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove all non-alphanumeric characters except dashes
	re := regexp.MustCompile(`[^\w-]`)
	s = re.ReplaceAllString(s, "")

	// Replace multiple dashes with a single dash
	re2 := regexp.MustCompile(`[-]+`)
	s = re2.ReplaceAllString(s, "-")

	// Trim leading and trailing dashes
	s = strings.Trim(s, "-")

	return s
}
