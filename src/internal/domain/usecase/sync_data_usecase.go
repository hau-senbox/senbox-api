package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
	"gorm.io/datatypes"
)

type SyncDataUsecase struct {
	SheetService   *sheets.Service
	SubmissionRepo *repository.SubmissionRepository
	SyncQueueRepo  *repository.SyncQueueRepository
}

type CreateFormAnswerRequest struct {
	SubmissionID    uint64
	SubmittedAt     time.Time
	StudentCustomID string
	UserCustomID    string
	FormCode        string
	FormName        string
	Answers         map[string]string
}

func (uc *SyncDataUsecase) CreateAndSyncFormAnswerv2(
	req CreateFormAnswerRequest,
	spreadsheetID string,
	sheetName string,
	headers []interface{},
	headerIndex map[string]int,
	queueID uint64,
) error {
	// Load Vietnam timezone
	loc, errT := time.LoadLocation("Asia/Ho_Chi_Minh")
	if errT != nil {
		return fmt.Errorf("failed to load timezone: %w", errT)
	}

	// Convert to Vietnam time
	submittedAtVN := req.SubmittedAt.In(loc)
	tFormatted := submittedAtVN.Format("2006-01-02 15:04:05")

	// Parse SubmittedAt
	//tFormatted := req.SubmittedAt.Format("2006-01-02 15:04:05")

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
		// Cập nhật trạng thái queue về FAIL
		_ = uc.SyncQueueRepo.UpdateStatusByID(queueID, value.SyncQueueStatusFailed)

		return fmt.Errorf("failed to append data: %w", err)
	}

	return nil
}

func (uc *SyncDataUsecase) GetData2Sync(afterCreatedAt time.Time, formNote []string) ([]CreateFormAnswerRequest, error) {
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
			SubmittedAt:     sub.CreatedAt,
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

func (uc *SyncDataUsecase) ExcuteCreateAndSyncFormAnswer(req request.SyncDataRequest) (string, error) {
	const defaultStartTime = "2025-08-01T00:00:00Z"

	// Marshal form notes để dùng truy vấn
	notesJSON, err := json.Marshal(req.FormNotes)
	if err != nil {
		return "", fmt.Errorf("failed to marshal form notes: %w", err)
	}

	// Truy vấn SyncQueue nếu đã có
	var (
		afterCreatedAt time.Time
		syncQueue      *entity.SyncQueue
	)

	existingQueue, err := uc.SyncQueueRepo.GetBySheetUrlAndSheetNameAndFormNotes(req.SheetUrl, req.SheetName, notesJSON)
	if err == nil && existingQueue != nil {
		// Nếu đã tồn tại bản ghi trước đó → dùng LastSubmittedAt
		afterCreatedAt = existingQueue.LastSubmittedAt
		if err != nil {
			return "", fmt.Errorf("invalid LastSubmittedAt in existing SyncQueue: %w", err)
		}
		syncQueue = existingQueue // sẽ dùng để update sau
	} else {
		// Không có dữ liệu cũ → bắt đầu từ ngày 1/6/2025
		afterCreatedAt, _ = time.Parse(time.RFC3339, defaultStartTime)
		syncQueue = &entity.SyncQueue{} // sẽ tạo mới sau
	}

	// Lấy spreadsheetID
	spreadsheetID, err := ExtractSpreadsheetID(req.SheetUrl)
	if err != nil {
		return "", fmt.Errorf("invalid sheet URL: %w", err)
	}

	// Lấy dữ liệu để đồng bộ
	dataList, err := uc.GetData2Sync(afterCreatedAt, req.FormNotes)
	if err != nil {
		return "", fmt.Errorf("failed to get data to sync: %w", err)
	}

	if len(dataList) == 0 {
		return "", errors.New("no data to sync")
	}

	// Chuẩn bị headers
	allAnswers := make([]map[string]string, len(dataList))
	for i, item := range dataList {
		allAnswers[i] = item.Answers
	}

	headers, headerIndex, err := uc.prepareHeaders(spreadsheetID, req.SheetName, allAnswers)
	if err != nil {
		return "", err
	}

	// Cập nhật thông tin cho SyncQueue (dù là tạo mới hay update)
	syncQueue.LastSubmissionID = dataList[len(dataList)-1].SubmissionID
	syncQueue.LastSubmittedAt = dataList[len(dataList)-1].SubmittedAt
	syncQueue.FormNotes = datatypes.JSON(notesJSON)
	syncQueue.SheetName = req.SheetName
	syncQueue.SpreadsheetID = spreadsheetID
	syncQueue.SheetUrl = req.SheetUrl
	syncQueue.Status = value.SyncQueueStatusPending
	syncQueue.IsAuto = req.IsAuto

	// Nếu là bản ghi cũ → update, nếu là bản mới → create
	if existingQueue != nil {
		if err := uc.SyncQueueRepo.Update(syncQueue); err != nil {
			return "", fmt.Errorf("failed to update sync queue: %w", err)
		}
	} else {
		if err := uc.SyncQueueRepo.Create(syncQueue); err != nil {
			return "", fmt.Errorf("failed to create sync queue: %w", err)
		}
	}

	// Đồng bộ dữ liệu lên Google Sheet ở nền
	go func(queueID uint64) {
		for i, item := range dataList {
			if err := uc.CreateAndSyncFormAnswerv2(item, spreadsheetID, req.SheetName, headers, headerIndex, queueID); err != nil {
				fmt.Printf("[SYNC ERROR] StudentCustomID %s: %v\n", item.StudentCustomID, err)
				continue
			}
			time.Sleep(1 * time.Second)
			if (i+1)%20 == 0 {
				fmt.Println("Reached 20 items, waiting 40 seconds...")
				time.Sleep(40 * time.Second)
			}
		}

		// Đánh dấu đã xong
		_ = uc.SyncQueueRepo.UpdateStatus(queueID, string(value.SyncQueueStatusDone))
	}(syncQueue.ID)

	return dataList[len(dataList)-1].SubmittedAt.String(), nil
}

func ExtractSpreadsheetID(sheetUrl string) (string, error) {
	re := regexp.MustCompile(`\/d\/([a-zA-Z0-9-_]+)`)
	matches := re.FindStringSubmatch(sheetUrl)
	if len(matches) < 2 {
		return "", fmt.Errorf("spreadsheet ID not found in URL")
	}
	return matches[1], nil
}

func (uc *SyncDataUsecase) prepareHeaders(spreadsheetID, sheetName string, allAnswers []map[string]string) ([]interface{}, map[string]int, error) {
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

func (uc *SyncDataUsecase) HasPendingSyncQueue() (bool, error) {
	ok, err := uc.SyncQueueRepo.HasPendingQueue()
	if err != nil {
		return false, fmt.Errorf("failed to check sync queue: %w", err)
	}
	if !ok {
		// Có queue đang pending → không thể tiếp tục
		return false, nil
	}
	// Không có pending → có thể sync
	return true, nil
}

func (uc *SyncDataUsecase) GetAllSyncQueue() ([]response.SyncQueueResponse, error) {
	queues, err := uc.SyncQueueRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var result []response.SyncQueueResponse
	for _, q := range queues {
		var formQRs []string

		// Giải mã JSON dạng []string từ cột FormNotes
		if err := json.Unmarshal(q.FormNotes, &formQRs); err != nil {
			formQRs = []string{}
		}

		result = append(result, response.SyncQueueResponse{
			ID:        q.ID,
			SheetURL:  q.SheetUrl,
			SheetName: q.SheetName,
			FormQRs:   strings.Join(formQRs, ","),
			IsAuto:    &q.IsAuto,
		})
	}

	return result, nil
}

func (uc *SyncDataUsecase) AutoSyncFormAnswersDaily() {
	queues, err := uc.SyncQueueRepo.GetAllAutoSync()
	if err != nil {
		log.Printf("Failed to fetch auto-sync queues: %v\n", err)
		return
	}

	for _, queue := range queues {
		var formNotesArr []string
		if err := json.Unmarshal(queue.FormNotes, &formNotesArr); err != nil {
			log.Printf("[AUTO SYNC ERROR] Failed to unmarshal FormNotes for QueueID %d: %v", queue.ID, err)
			continue
		}

		req := request.SyncDataRequest{
			SheetUrl:  queue.SheetUrl,
			SheetName: queue.SheetName,
			FormNotes: formNotesArr,
		}

		go func(q entity.SyncQueue, r request.SyncDataRequest) {
			log.Printf("[AUTO SYNC] Start syncing for Sheet: %s", q.SheetName)
			_, err := uc.ExcuteCreateAndSyncFormAnswer(r)
			if err != nil {
				log.Printf("[AUTO SYNC ERROR] QueueID %d: %v", q.ID, err)
			}
		}(queue, req)

		// doi 5 phut
		time.Sleep(5 * time.Minute)
	}

}

func (uc *SyncDataUsecase) StartAutoSyncScheduler() {
	c := cron.New(cron.WithSeconds())
	// chay vao lic 00:00
	_, err := c.AddFunc("0 0 0 * * *", func() {
		log.Println("[CRON] Running AutoSyncFormAnswersDaily at", time.Now().Format(time.RFC3339))
		uc.AutoSyncFormAnswersDaily()
	})

	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	// Chạy mỗi giây
	// _, err = c.AddFunc("@every 1s", func() {
	// 	log.Println("[CRON] Running job every second at", time.Now().Format(time.RFC3339))
	// 	uc.AutoSyncFormAnswersDaily()
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to add every-second job: %v", err)
	// }

	c.Start()
}
