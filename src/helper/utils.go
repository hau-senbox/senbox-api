package helper

import (
	"regexp"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"strings"
	"time"
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

func ParseAtrValueStringToStruct(s string) request.AtrValueString {
	result := request.AtrValueString{}
	pairs := strings.Split(s, ";")

	for _, pair := range pairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])

		switch key {
		case "question_key":
			result.QuestionKey = &valueStr
		case "question_db":
			result.QuestionDB = &valueStr
		case "time_sort":
			result.TimeSort = value.TimeSort(valueStr)
		case "user_id":
			result.UserID = valueStr
		case "duration":
			dates := strings.Split(valueStr, ",")
			if len(dates) == 2 {
				start, err1 := parseDate(dates[0])
				end, err2 := parseDate(dates[1])
				if err1 == nil && err2 == nil {
					result.Duration = &value.TimeRange{
						Start: start,
						End:   end,
					}
				}
			}
		}

	}

	return result
}

func parseDate(s string) (time.Time, error) {
	// "2/3/2025-00:00" → layout: "2/1/2006-15:04"
	return time.Parse("2/1/2006-15:04", strings.TrimSpace(s))
}
