package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
		Visible bool   `json:"visible"`
		Icon    string `json:"icon"`
		Color   string `json:"color"`
		URL     string `json:"url"`
		FormQR  string `json:"form_qr"`
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
		newVal["form_qr"] = old.FormQR
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

func WriteDataToSheet(spreadsheetID, sheetName, startCell string, values [][]interface{}, credentialsPath string) error {
	ctx := context.Background()

	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %w", err)
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("failed to parse credentials: %w", err)
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create sheets service: %w", err)
	}

	writeRange := fmt.Sprintf("%s!%s", sheetName, startCell)
	valueRange := &sheets.ValueRange{Values: values}

	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, valueRange).
		ValueInputOption("USER_ENTERED").
		Do()

	if err != nil {
		return fmt.Errorf("failed to write data to sheet: %w", err)
	}

	return nil
}

func AppendDataToSheet(spreadsheetID, sheetName string, values [][]interface{}, credentialsPath string) error {
	ctx := context.Background()

	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return fmt.Errorf("unable to read credentials file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("unable to parse credentials file to config: %v", err)
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Gọi append, dùng USER_ENTERED để tự động xử lý định dạng (ví dụ ngày/tháng)
	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, sheetName, valueRange).
		ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		Do()

	if err != nil {
		return fmt.Errorf("unable to append data to sheet: %v", err)
	}

	return nil
}

func GetSheetHeaders(spreadsheetID, sheetName, credentialsPath string) ([]string, error) {
	ctx := context.Background()

	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	readRange := fmt.Sprintf("%s!1:1", sheetName)
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	headers := []string{}
	if len(resp.Values) > 0 {
		for _, cell := range resp.Values[0] {
			if str, ok := cell.(string); ok {
				headers = append(headers, str)
			} else {
				headers = append(headers, "")
			}
		}
	}

	return headers, nil
}

func SheetHasData(spreadsheetID, sheetName, credentialsPath string) (bool, error) {
	ctx := context.Background()

	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return false, err
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return false, err
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return false, err
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, sheetName+"!A1:A2").Do()
	if err != nil {
		return false, err
	}

	// Nếu có ít nhất 1 dòng dữ liệu
	return len(resp.Values) > 0, nil
}

func GetSheetHeader(spreadsheetID, sheetName, credentialsPath string) ([]string, error) {
	srv, err := GetSheetsService(credentialsPath)
	if err != nil {
		return nil, err
	}

	readRange := sheetName + "!1:1"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Values) == 0 {
		return []string{}, nil
	}

	headers := make([]string, len(resp.Values[0]))
	for i, val := range resp.Values[0] {
		headers[i] = fmt.Sprintf("%v", val)
	}
	return headers, nil
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

func AppendFormAnswersToSheet(srv *sheets.Service, spreadsheetID, sheetName string, baseInfo []interface{}, answers map[string]string) error {
	readRange := fmt.Sprintf("%s!1:1", sheetName)
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return fmt.Errorf("failed to read header row: %w", err)
	}

	// Bước 1: Khởi tạo header (dòng đầu tiên)
	var headers []interface{}
	headerIndex := make(map[string]int)

	if len(resp.Values) == 0 || len(resp.Values[0]) == 0 {
		// Sheet chưa có header → thêm mặc định 5 cột đầu
		defaultHeaders := []string{"SubmittedAt", "StudentID", "UserID", "FormCode", "FormName"}
		for _, h := range defaultHeaders {
			headers = append(headers, h)
		}
	} else {
		headers = resp.Values[0]
	}

	// Bước 2: Mapping header index
	for idx, h := range headers {
		headerIndex[fmt.Sprintf("%v", h)] = idx
	}

	// Bước 3: Mở rộng header nếu có câu hỏi mới
	for q := range answers {
		if _, ok := headerIndex[q]; !ok {
			headers = append(headers, q)
			headerIndex[q] = len(headers) - 1
		}
	}

	// Bước 4: Nếu có thêm header mới hoặc chưa có dòng đầu tiên → update lại dòng header
	if len(resp.Values) == 0 || len(headers) > len(resp.Values[0]) {
		updateRange := fmt.Sprintf("%s!1:1", sheetName)
		_, err := srv.Spreadsheets.Values.Update(spreadsheetID, updateRange, &sheets.ValueRange{
			Values: [][]interface{}{headers},
		}).ValueInputOption("RAW").Do()
		if err != nil {
			return fmt.Errorf("failed to update headers: %w", err)
		}
	}

	// Bước 5: Tạo dòng dữ liệu mới theo đúng thứ tự header
	row := make([]interface{}, len(headers))
	copy(row, baseInfo) // baseInfo = [SubmittedAt, StudentID, UserID, FormCode, FormName]

	for q, a := range answers {
		if colIndex, ok := headerIndex[q]; ok {
			row[colIndex] = a
		}
	}

	// Bước 6: Ghi vào sheet
	appendRange := fmt.Sprintf("%s!A:Z", sheetName)
	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, appendRange, &sheets.ValueRange{
		Values: [][]interface{}{row},
	}).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do()
	if err != nil {
		return fmt.Errorf("failed to append data: %w", err)
	}

	return nil
}
