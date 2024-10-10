package repository

import (
	"encoding/json"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"
	"time"

	"gorm.io/gorm"
)

type SubmissionRepository struct {
	DBConn *gorm.DB
}

type Messaging struct {
	Email        []string `json:"email" binding:"required"`
	Value3       []string `json:"value3" binding:"required"`
	MessageBox   *string  `json:"messageBox"`
	QuestionType string   `json:"questionType" binding:"required"`
}

type SubmissionDataItem struct {
	QuestionId string     `json:"question_id" binding:"required"`
	Question   string     `json:"question" binding:"required"`
	Answer     string     `json:"answer" binding:"required"`
	Messaging  *Messaging `json:"messaging"`
}

type SubmissionData struct {
	Items []SubmissionDataItem `json:"items" binding:"required"`
}
type CreateSubmissionParams struct {
	FormId             uint64
	FormName           string
	FormNote           string
	FormSpreadsheetUrl string
	DeviceId           string
	DeviceFirstValue   string
	DeviceSecondValue  string
	DeviceThirdValue   string
	DeviceName         string
	DeviceNote         string
	SpreadsheetId      string
	SheetName          string
	SubmissionData     SubmissionData
	SubmissionType     value.SubmissionType
	OpenedAt           time.Time
}

func (receiver *SubmissionRepository) CreateSubmission(params CreateSubmissionParams) error {
	items := make([]entity.SubmissionDataItem, 0)
	for _, item := range params.SubmissionData.Items {
		var msg *entity.Messaging = nil
		if item.Messaging != nil {
			msg = &entity.Messaging{
				Email:        item.Messaging.Email,
				Value3:       item.Messaging.Value3,
				MessageBox:   item.Messaging.MessageBox,
				QuestionType: item.Messaging.QuestionType,
			}
		}
		items = append(items, entity.SubmissionDataItem{
			QuestionId: item.QuestionId,
			Question:   item.Question,
			Answer:     item.Answer,
			Messaging:  msg,
		})
	}

	data := entity.SubmissionData{
		Items: items,
	}

	dataInJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	submission := entity.SSubmission{
		FormId:             params.FormId,
		FormName:           params.FormName,
		FormNote:           params.FormNote,
		FormSpreadsheetUrl: params.FormSpreadsheetUrl,
		DeviceId:           params.DeviceId,
		DeviceFirstValue:   params.DeviceFirstValue,
		DeviceSecondValue:  params.DeviceSecondValue,
		DeviceThirdValue:   params.DeviceThirdValue,
		DeviceName:         params.DeviceName,
		DeviceNote:         params.DeviceNote,
		SpreadsheetId:      params.SpreadsheetId,
		SheetName:          params.SheetName,
		SubmissionData:     dataInJSON,
		SubmissionType:     params.SubmissionType,
		Status:             value.SubmissionStatusAccepted,
		OpenedAt:           params.OpenedAt,
	}

	return receiver.DBConn.Create(&submission).Error
}

func (receiver *SubmissionRepository) FindFirstPendingSync() (entity.SSubmission, error) {
	var submission entity.SSubmission
	err := receiver.DBConn.
		Order("status asc").
		First(&submission, "status in ? AND number_attempt < ?", []value.SubmissionStatus{value.SubmissionStatusAccepted, value.SubmissionStatusAttempted}, 10).Error
	if err != nil {
		return submission, err
	}

	return submission, nil
}

func (receiver *SubmissionRepository) MarkStatusSucceeded(id uint64) error {
	return receiver.DBConn.Model(&entity.SSubmission{}).
		Where("id = ?", id).
		Update("status", value.SubmissionStatusSynced).Error
}

func (receiver *SubmissionRepository) MarkStatusAttempted(id uint64) error {
	return receiver.DBConn.
		Exec("UPDATE s_submission SET status = ?, number_attempt = number_attempt + 1 WHERE id = ?", value.SubmissionStatusAttempted, id).Error
}

func (receiver *SubmissionRepository) FindRecentByFormId(formId uint64) (entity.SSubmission, error) {
	var submission entity.SSubmission
	err := receiver.DBConn.Where("form_id = ?", formId).Order("created_at DESC").First(&submission).Error
	if err != nil {
		return entity.SSubmission{}, err
	}

	return submission, err
}

func (receiver *SubmissionRepository) DublicateSubmissions(params CreateSubmissionParams) error {
	items := make([]entity.SubmissionDataItem, 0)
	for _, item := range params.SubmissionData.Items {
		var msg *entity.Messaging = nil
		if item.Messaging != nil {
			msg = &entity.Messaging{
				Email:        item.Messaging.Email,
				Value3:       item.Messaging.Value3,
				MessageBox:   item.Messaging.MessageBox,
				QuestionType: item.Messaging.QuestionType,
			}
		}
		items = append(items, entity.SubmissionDataItem{
			QuestionId: item.QuestionId,
			Question:   item.Question,
			Answer:     item.Answer,
			Messaging:  msg,
		})
	}

	data := entity.SubmissionData{
		Items: items,
	}

	dataInJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	submission := entity.SSubmission{
		FormId:             params.FormId,
		FormName:           params.FormName,
		FormNote:           params.FormNote,
		FormSpreadsheetUrl: params.FormSpreadsheetUrl,
		DeviceId:           params.DeviceId,
		DeviceFirstValue:   params.DeviceFirstValue,
		DeviceSecondValue:  params.DeviceSecondValue,
		DeviceThirdValue:   params.DeviceThirdValue,
		DeviceName:         params.DeviceName,
		DeviceNote:         params.DeviceNote,
		SpreadsheetId:      params.SpreadsheetId,
		SheetName:          params.SheetName,
		SubmissionData:     dataInJSON,
		SubmissionType:     params.SubmissionType,
		Status:             value.SubmissionStatusAccepted,
		OpenedAt:           params.OpenedAt,
	}

	return receiver.DBConn.Create(&submission).Error
}

func (receiver *SubmissionRepository) FindFirstPendingPrioritizedSync() (entity.SSubmission, error) {
	var submission entity.SSubmission
	err := receiver.DBConn.InnerJoins("INNER JOIN s_form on s_submission.form_id = s_form.form_id").
		Order("status asc").
		First(&submission, "s_submission.status in ? AND number_attempt < ? AND s_form.sync_strategy = ?",
			[]value.SubmissionStatus{value.SubmissionStatusAccepted, value.SubmissionStatusAttempted},
			10,
			value.FormSyncStrategyOnSubmit).
		Error
	if err != nil {
		return submission, err
	}

	return submission, nil
}
