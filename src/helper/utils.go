package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"sort"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
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
		Visible      bool   `json:"visible"`
		Icon         string `json:"icon"`
		Color        string `json:"color"`
		URL          string `json:"url"`
		FormQR       string `json:"form_qr"`
		ShowedTop    string `json:"showed_top"`
		ShowedBottom string `json:"showed_bottom"`
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
		"color":         old.Color,
		"icon":          old.Icon,
		"visible":       old.Visible,
		"showed_top":    old.ShowedTop,
		"showed_bottom": old.ShowedBottom,
	}
	if isButtonForm {
		newVal["form_qr"] = old.FormQR
	} else {
		newVal["url"] = old.URL
	}

	// Build object ngoài
	wrapped := map[string]interface{}{
		"id":            comp.ID.String(),
		"name":          comp.Name,
		"type":          string(comp.Type),
		"key":           comp.Key,
		"color":         old.Color,
		"icon":          old.Icon,
		"visible":       old.Visible,
		"value":         newVal,
		"showed_top":    old.ShowedTop,
		"showed_bottom": old.ShowedBottom,
	}

	// field chính ở ngoài: form_qr hoặc url
	if isButtonForm {
		wrapped["form_qr"] = old.FormQR
	} else {
		wrapped["url"] = old.URL
	}

	jsonBytes, err := json.Marshal(wrapped)
	if err != nil {
		return oldValue
	}

	return string(jsonBytes)
}

func FormatProfileLink(rawURL string, id string) string {
	parts := strings.Split(rawURL, "/")
	if len(parts) > 1 {
		lastPart := parts[len(parts)-1]
		parts[len(parts)-1] = "[" + lastPart + "]:" + id
		return strings.Join(parts, "/")
	}
	return rawURL
}

func FilterUsersByName(users []response.UserResponse, name string) []response.UserResponse {
	name = strings.ToLower(name)
	result := make([]response.UserResponse, 0)
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Nickname), name) {
			result = append(result, u)
			sort.Slice(result, func(i, j int) bool {
				return strings.ToLower(result[i].Nickname) < strings.ToLower(result[j].Nickname)
			})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].Nickname) < strings.ToLower(result[j].Nickname)
	})
	return result
}

func FilterChildrenByName(users []response.ChildrenResponse, name string) []response.ChildrenResponse {
	name = strings.ToLower(name)
	result := make([]response.ChildrenResponse, 0)
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.ChildName), name) {
			result = append(result, u)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].ChildName) < strings.ToLower(result[j].ChildName)
	})
	return result
}

func FilterStudentByName(users []response.StudentResponse, name string) []response.StudentResponse {
	name = strings.ToLower(name)
	result := make([]response.StudentResponse, 0)
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.StudentName), name) {
			result = append(result, u)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].StudentName) < strings.ToLower(result[j].StudentName)
	})
	return result
}

func FilterTeacherByName(users []response.TeacherResponse, name string) []response.TeacherResponse {
	name = strings.ToLower(name)
	result := make([]response.TeacherResponse, 0)
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.TeacherName), name) {
			result = append(result, u)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].TeacherName) < strings.ToLower(result[j].TeacherName)
	})
	return result
}

func FilterStaffByName(users []response.StaffResponse, name string) []response.StaffResponse {
	name = strings.ToLower(name)
	result := make([]response.StaffResponse, 0)
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.StaffName), name) {
			result = append(result, u)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].StaffName) < strings.ToLower(result[j].StaffName)
	})
	return result
}

func GetSheetsService(credentialsPath string) (*sheets.Service, error) {
	ctx := context.Background()

	// Đọc file credentials JSON
	b, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}

	// Parse credentials
	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	// Tạo sheets service
	srv, err := sheets.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx)))
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func ParseFlexibleTime(value string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05.000 -0700 MST",
		"2006-01-02 15:04:05.00 -0700 MST",
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05", // no offset, no zone
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t.UTC(), nil
		}
	}

	// Fallback: parse as if it's in UTC
	layout := "2006-01-02 15:04:05.00"
	return time.ParseInLocation(layout, value, time.UTC)
}
