package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/model"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/messaging"
	"sen-global-api/pkg/sheet"

	firebase "firebase.google.com/go/v4"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"gorm.io/gorm"
)

type SubmitFormUseCase struct {
	*repository.DeviceRepository
	*repository.UserEntityRepository
	*repository.FormRepository
	*repository.QuestionRepository
	*repository.SubmissionRepository
	*repository.SettingRepository
	*repository.MobileDeviceRepository
	*repository.FormQuestionRepository
	*repository.CodeCountingRepository
	*repository.AnswerRepository
	*sheet.Writer
	*sheet.Reader
	DriveService        *drive.Service
	OutputSpreadsheetID string
	FirebaseApp         *firebase.App
	DB                  *gorm.DB
}

func (receiver *SubmitFormUseCase) AnswerForm(id uint64, req request.SubmitFormRequest) error {
	form, err := receiver.GetFormByID(id)
	if err != nil {
		return err
	}

	return receiver.answerFormSaveToFormOutputSheet(form, req)
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func (receiver *SubmitFormUseCase) answerFormSaveToFormOutputSheet(form *entity.SForm, req request.SubmitFormRequest) error {
	submissionItems := make([]repository.SubmissionDataItem, 0)
	questions, err := receiver.GetQuestionsByIDs(Map(req.Answers, func(answer request.Answer) string { return answer.QuestionID }))
	if err != nil {
		return fmt.Errorf("system cannot find questions for this form: %s", form.Name)
	}

	rememberAnswers := make([]entity.MemoryComponentValue, 0)
	for _, answer := range req.Answers {
		for _, question := range questions {
			if answer.QuestionID == question.ID.String() {
				if answer.Remember {
					rememberAnswers = append(rememberAnswers, entity.MemoryComponentValue{
						ComponentName: question.QuestionType,
						Value:         answer.Answer,
					})
				}

				submissionItems = append(submissionItems, repository.SubmissionDataItem{
					QuestionID: question.ID.String(),
					Key:        answer.Key,
					DB:         answer.DB,
					Question:   question.Question,
					Answer:     answer.Answer,
				})
			}
		}
	}

	submissionData := repository.SubmissionData{
		Items: submissionItems,
	}
	createSubmissionParams := repository.CreateSubmissionParams{
		FormID:         form.ID,
		UserID:         req.UserID,
		SubmissionData: submissionData,
		OpenedAt:       req.OpenedAt,
		CustomID:       req.CustomID,
	}
	submissionID, err := receiver.CreateSubmission(createSubmissionParams)
	if err != nil {
		log.Error("SubmitFormUseCase.answerFormSaveToFormOutputSheet", err)
		return errors.New("system cannot handle the submission")
	}

	if len(rememberAnswers) > 0 {
		err = receiver.QuestionRepository.CreateMemoryComponentValuesDuplicate(rememberAnswers)
		if err != nil {
			log.Error("SubmitFormUseCase.answerFormSaveToFormOutputSheet", err)
			return errors.New("system cannot handle the submission memory component values")
		}
	}

	// tao anwser
	for _, item := range submissionItems {

		if item.Key == "" && item.DB == "" {
			continue
		}

		ansJSON, _ := json.Marshal(item.Answer)

		answer := &entity.SAnswer{
			ID:           uuid.New(),
			UserID:       req.UserID,
			SubmissionID: submissionID,
			Key:          item.Key,
			DB:           item.DB,
			Response:     json.RawMessage(ansJSON),
		}

		err := receiver.AnswerRepository.Create(answer)
		if err != nil {
			log.Error("SubmitFormUseCase.answerFormSaveToFormOutputSheet: create answer", err)
		}
	}

	defer func() {
		receiver.sendNotification(form)
	}()

	return nil
}

func (receiver *SubmitFormUseCase) sendNotification(form *entity.SForm) {
	questions, err := receiver.GetQuestionsByFormID(form.ID)
	if err != nil {
		log.Error("Form ", form.Note, " has not send notification question")
		return
	}
	sendNotificationQuestion := model.FormQuestionItem{}
	for _, q := range questions {
		if q.QuestionType == value.GetStringValue(value.QuestionSendNotification) {
			sendNotificationQuestion = q
		}
	}
	if sendNotificationQuestion.ID == "" {
		log.Debug("Form ", form.Note, " has not send notification question")
		return
	}
	type QuestionAttributes struct {
		Value string `json:"value"`
	}
	var att QuestionAttributes
	err = json.Unmarshal(sendNotificationQuestion.Attributes, &att)
	if err != nil {
		log.Error("Can not unmarshal send notification value ", sendNotificationQuestion.Attributes)
		return
	}

	md, err := receiver.FindByDeviceID(att.Value, receiver.DB)
	if err != nil {
		log.Error("FCM Token could not be found for the device id ", att.Value)
	}

	noti := messaging.NotificationParams{
		Title:       "New Form Submit",
		Message:     "",
		DeviceToken: md.FCMToken,
		Type:        value.NotificationType_NewFormSubmit,
	}
	err = messaging.SendNotification(receiver.FirebaseApp, noti)
	if err != nil {
		log.Error("Failed to send notification ", err)
	}
}
