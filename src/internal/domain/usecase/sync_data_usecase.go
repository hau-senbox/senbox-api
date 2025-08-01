package usecase

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"time"

	"google.golang.org/api/sheets/v4"
	"gorm.io/datatypes"
)

type SyncDataUsecae struct {
	SheetService   *sheets.Service
	SubmissionRepo *repository.SubmissionRepository
	SyncQueueRepo  *repository.SyncQueueRepository
}

type CreateFormAnswerRequest struct {
	SubmissionID    uint64
	SubmittedAt     string
	StudentCustomID string
	UserCustomID    string
	FormCode        string
	FormName        string
	Answers         map[string]string
}

func (uc *SyncDataUsecae) CreateAndSyncFormAnswer(req CreateFormAnswerRequest, spreadsheetID string, sheetName string) error {
	// Parse SubmittedAt
	tFormatted := req.SubmittedAt
	if t, err := time.Parse("2006-01-02 15:04:05.999 -0700 MST", req.SubmittedAt); err == nil {
		tFormatted = t.Format("2006-01-02 15:04:05")
	}

	// Build base info
	baseInfo := []interface{}{
		tFormatted,
		req.StudentCustomID,
		req.UserCustomID,
		req.FormCode,
		req.FormName,
	}

	// Đọc dòng header hiện tại
	readRange := fmt.Sprintf("%s!1:1", sheetName)
	resp, err := uc.SheetService.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return fmt.Errorf("failed to read sheet headers: %w", err)
	}

	var headers []interface{}
	headerIndex := make(map[string]int)

	// Nếu chưa có header → khởi tạo
	if len(resp.Values) == 0 || len(resp.Values[0]) == 0 {
		defaultHeaders := []string{"SubmittedAt", "StudentCustomID", "UserCustomID", "FormCode", "FormName"}
		for _, h := range defaultHeaders {
			headers = append(headers, h)
		}
	} else {
		headers = resp.Values[0]
	}

	// Mapping header -> index
	for idx, h := range headers {
		headerIndex[fmt.Sprintf("%v", h)] = idx
	}

	// Bổ sung header nếu có câu hỏi mới
	for q := range req.Answers {
		if _, ok := headerIndex[q]; !ok {
			headers = append(headers, q)
			headerIndex[q] = len(headers) - 1
		}
	}

	// Cập nhật lại header nếu có thay đổi
	if len(resp.Values) == 0 || len(headers) > len(resp.Values[0]) {
		updateRange := fmt.Sprintf("%s!1:1", sheetName)
		_, err := uc.SheetService.Spreadsheets.Values.Update(spreadsheetID, updateRange, &sheets.ValueRange{
			Values: [][]interface{}{headers},
		}).ValueInputOption("RAW").Do()
		time.Sleep(2 * time.Second)
		if err != nil {
			return fmt.Errorf("failed to update headers: %w", err)
		}
	}

	// Tạo dòng dữ liệu mới
	row := make([]interface{}, len(headers))
	copy(row, baseInfo) // Gán SubmittedAt, StudentID, UserID...

	for q, a := range req.Answers {
		if colIndex, ok := headerIndex[q]; ok {
			row[colIndex] = a
		}
	}

	// Ghi dòng mới vào sheet
	appendRange := fmt.Sprintf("%s!A:Z", sheetName)
	_, err = uc.SheetService.Spreadsheets.Values.Append(spreadsheetID, appendRange, &sheets.ValueRange{
		Values: [][]interface{}{row},
	}).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do()
	if err != nil {
		return fmt.Errorf("failed to append data: %w", err)
	}

	return nil
}

func (uc *SyncDataUsecae) CreateAndSyncFormAnswerv2(
	req CreateFormAnswerRequest,
	spreadsheetID string,
	sheetName string,
	headers []interface{},
	headerIndex map[string]int,
) error {
	// Parse SubmittedAt
	tFormatted := req.SubmittedAt
	if t, err := time.Parse("2006-01-02 15:04:05.999 -0700 MST", req.SubmittedAt); err == nil {
		tFormatted = t.Format("2006-01-02 15:04:05")
	}

	// Base info
	baseInfo := []interface{}{
		tFormatted,
		req.StudentCustomID,
		req.UserCustomID,
		req.FormCode,
		req.FormName,
	}

	// Tạo dòng dữ liệu mới
	row := make([]interface{}, len(headers))
	copy(row, baseInfo)

	for q, a := range req.Answers {
		if colIndex, ok := headerIndex[q]; ok {
			row[colIndex] = a
		}
	}

	appendRange := fmt.Sprintf("%s!A:Z", sheetName)
	_, err := uc.SheetService.Spreadsheets.Values.Append(spreadsheetID, appendRange, &sheets.ValueRange{
		Values: [][]interface{}{row},
	}).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do()
	if err != nil {
		return fmt.Errorf("failed to append data: %w", err)
	}

	return nil
}

func (uc *SyncDataUsecae) GetData2Sync(afterCreatedAt time.Time, formNote []string) ([]CreateFormAnswerRequest, error) {
	// 1. Lấy danh sách submission mới nhất
	submissions, err := uc.SubmissionRepo.GetSubmissionByCreatedAtAndForms(afterCreatedAt, formNote)
	if err != nil {
		return nil, fmt.Errorf("failed to get submissions: %w", err)
	}

	var results []CreateFormAnswerRequest

	for _, sub := range submissions {
		var submissionData entity.SubmissionData
		if err := json.Unmarshal(sub.SubmissionData, &submissionData); err != nil {
			continue // Bỏ qua nếu parse lỗi
		}

		answers := make(map[string]string)
		for _, item := range submissionData.Items {
			answers[item.Question] = item.Answer
		}

		req := CreateFormAnswerRequest{
			SubmissionID:    sub.ID,
			SubmittedAt:     sub.CreatedAt.String(),
			StudentCustomID: sub.StudentCustomID,
			UserCustomID:    sub.UserCustomID,
			FormCode:        sub.Form.Note,
			FormName:        sub.Form.Name,
			Answers:         answers,
		}

		results = append(results, req)
	}

	return results, nil
}

func (uc *SyncDataUsecae) ExcuteCreateAndSyncFormAnswer(req request.SyncDataRequest) (string, error) {
	// Parse thời gian từ chuỗi
	afterCreatedAt, err := time.Parse(time.RFC3339, req.LastSubmitTime)
	if err != nil {
		return "", fmt.Errorf("invalid time format (must be RFC3339): %w", err)
	}

	// Lấy ID từ sheet URL
	spreadsheetID, err := ExtractSpreadsheetID(req.SheetUrl)
	if err != nil {
		return "", fmt.Errorf("invalid sheet URL: %w", err)
	}

	// Lấy dữ liệu cần đồng bộ
	dataList, err := uc.GetData2Sync(afterCreatedAt, req.FormNotes)
	if err != nil {
		return "", fmt.Errorf("failed to get data to sync: %w", err)
	}

	// Nếu không có gì để sync thì trả về luôn
	if len(dataList) == 0 {
		return "", nil
	}

	// Lấy submissionID/timestamp cuối cùng
	latestSubmissionTime := dataList[len(dataList)-1].SubmittedAt
	// Collect all answers
	allAnswers := make([]map[string]string, len(dataList))
	for i, item := range dataList {
		allAnswers[i] = item.Answers
	}

	// Chỉ đọc header 1 lần
	headers, headerIndex, err := uc.prepareHeaders(spreadsheetID, req.SheetName, allAnswers)
	if err != nil {
		return "", err
	}
	// Parse SubmittedAt từ string → time.Time
	layout := "2006-01-02 15:04:05.000 -0700 -07"
	parsedSubmittedAt, err := time.Parse(layout, dataList[len(dataList)-1].SubmittedAt)

	if err != nil {
		return "", fmt.Errorf("invalid SubmittedAt format: %w", err)
	}

	// Marshal FormNotes
	notesJSON, err := json.Marshal(req.FormNotes)
	if err != nil {
		return "", fmt.Errorf("failed to marshal form notes: %w", err)
	}

	// Tạo SyncQueue trước khi chạy goroutine
	syncQueue := &entity.SyncQueue{
		LastSubmissionID: dataList[len(dataList)-1].SubmissionID,
		LastSubmittedAt:  parsedSubmittedAt,
		FormNotes:        datatypes.JSON(notesJSON),
		SheetName:        req.SheetName,
		SpreadsheetID:    spreadsheetID,
		Status:           value.SyncQueueStatusPending, // pending
	}

	// Lưu queue vào DB
	if err := uc.SyncQueueRepo.Create(syncQueue); err != nil {
		return "", fmt.Errorf("failed to create sync queue: %w", err)
	}

	// Ghi từng dòng
	go func() {

		for _, item := range dataList {

			if err := uc.CreateAndSyncFormAnswerv2(item, spreadsheetID, req.SheetName, headers, headerIndex); err != nil {
				fmt.Printf("[SYNC ERROR] StudentCustomID %s: %v\n", item.StudentCustomID, err)
				continue
			}
			time.Sleep(1 * time.Second)
		}
		// Cập nhật queue: done
		_ = uc.SyncQueueRepo.UpdateStatus(syncQueue.ID, string(value.SyncQueueStatusDone))
	}()

	// Trả kết quả ngay lập tức (timestamp cuối cùng)
	return latestSubmissionTime, nil
}

func ExtractSpreadsheetID(sheetUrl string) (string, error) {
	re := regexp.MustCompile(`\/d\/([a-zA-Z0-9-_]+)`)
	matches := re.FindStringSubmatch(sheetUrl)
	if len(matches) < 2 {
		return "", fmt.Errorf("spreadsheet ID not found in URL")
	}
	return matches[1], nil
}

func (uc *SyncDataUsecae) prepareHeaders(spreadsheetID, sheetName string, allAnswers []map[string]string) ([]interface{}, map[string]int, error) {
	readRange := fmt.Sprintf("%s!1:1", sheetName)
	resp, err := uc.SheetService.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read sheet headers: %w", err)
	}

	var headers []interface{}
	headerIndex := make(map[string]int)

	if len(resp.Values) == 0 || len(resp.Values[0]) == 0 {
		headers = []interface{}{"SubmittedAt", "StudentCustomID", "UserCustomID", "FormCode", "FormName"}
	} else {
		headers = resp.Values[0]
	}

	// Map existing headers
	for idx, h := range headers {
		headerIndex[fmt.Sprintf("%v", h)] = idx
	}

	// Append new headers from all answers
	for _, ans := range allAnswers {
		for q := range ans {
			if _, ok := headerIndex[q]; !ok {
				headers = append(headers, q)
				headerIndex[q] = len(headers) - 1
			}
		}
	}

	// Update header if needed
	if len(resp.Values) == 0 || len(headers) > len(resp.Values[0]) {
		updateRange := fmt.Sprintf("%s!1:1", sheetName)
		_, err := uc.SheetService.Spreadsheets.Values.Update(spreadsheetID, updateRange, &sheets.ValueRange{
			Values: [][]interface{}{headers},
		}).ValueInputOption("RAW").Do()
		time.Sleep(2 * time.Second) // tránh rate limit tiếp
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update headers: %w", err)
		}
	}

	return headers, headerIndex, nil
}
