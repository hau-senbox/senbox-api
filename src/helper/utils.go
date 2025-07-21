package helper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sen-global-api/internal/domain/entity/components"
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
		case "key":
			result.Key = &valueStr
		case "db":
			result.DB = &valueStr
		case "sort":
			result.TimeSort = value.TimeSort(valueStr)
		case "user_id":
			result.UserID = valueStr
		case "date_duration":
			dates := strings.Split(valueStr, ",")
			if len(dates) == 2 {
				start, err1 := parseDate(dates[0])
				end, err2 := parseDate(dates[1])
				if err1 == nil && err2 == nil {
					result.DateDuration = &value.TimeRange{
						Start: start,
						End:   end,
					}
				}
			}
		case "quantity":
			result.Quantity = &valueStr
		}

	}

	return result
}

func parseDate(s string) (time.Time, error) {
	// "2/3/2025-00:00" → layout: "2/1/2006-15:04"
	return time.Parse("2/1/2006-15:04", strings.TrimSpace(s))
}

// func ParseAtrValueListStringToStructs(s string, userID string) []request.AtrValueString {
// 	// Làm sạch chuỗi: remove đầu/cuối []
// 	s = strings.TrimPrefix(s, "['")
// 	s = strings.TrimSuffix(s, "']")
// 	items := strings.Split(s, "','")

// 	var results []request.AtrValueString

// 	for _, item := range items {
// 		parsed := ParseAtrValueStringToStruct(item)
// 		parsed.UserID = userID // inject userID từ context
// 		if parsed.Quantity == 0 {
// 			parsed.Quantity = 1 // fallback nếu không có hoặc lỗi
// 		}
// 		results = append(results, parsed)
// 	}

// 	return results
// }

func GetVisibleToValueComponent(value string) (bool, error) {
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(value), &parsed); err != nil {
		return true, fmt.Errorf("failed to parse component value JSON: %w", err)
	}

	// Default visible to true
	visible := true
	if v, ok := parsed["visible"].(bool); ok {
		visible = v
	}
	return visible, nil
}

func BuildSectionValueMenu(oldValue string, comp components.Component) string {
	var old struct {
		Visible bool   `json:"visible"`
		Icon    string `json:"icon"`
		Color   string `json:"color"`
		URL     string `json:"url"`
	}

	err := json.Unmarshal([]byte(oldValue), &old)
	if err != nil {
		// fallback nếu lỗi unmarshal
		return oldValue
	}

	// Xác định field chính là "form_qr" hay "url"
	isButtonForm := comp.Type == "button_form"

	// Build value nội
	newVal := map[string]interface{}{
		"color":   old.Color,
		"icon":    old.Icon,
		"visible": old.Visible,
	}
	if isButtonForm {
		newVal["form_qr"] = old.URL
	} else {
		newVal["url"] = old.URL
	}

	// Build object ngoài
	wrapped := map[string]interface{}{
		"id":      comp.ID.String(),
		"name":    comp.Name,
		"type":    string(comp.Type),
		"key":     comp.Key,
		"color":   old.Color,
		"icon":    old.Icon,
		"visible": old.Visible,
		"value":   newVal,
	}

	// field chính ở ngoài: form_qr hoặc url
	if isButtonForm {
		wrapped["form_qr"] = old.URL
	} else {
		wrapped["url"] = old.URL
	}

	jsonBytes, err := json.Marshal(wrapped)
	if err != nil {
		return oldValue
	}

	return string(jsonBytes)
}
