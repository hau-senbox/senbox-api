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
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"

	firebase "firebase.google.com/go/v4"
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
	*sheet.Writer
	*sheet.Reader
	DriveService        *drive.Service
	OutputSpreadsheetId string
	FirebaseApp         *firebase.App
	DB                  *gorm.DB
}

func (receiver *SubmitFormUseCase) AnswerForm(id uint64, req request.SubmitFormRequest) error {
	form, err := receiver.GetFormById(id)
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
	questions, err := receiver.GetQuestionsByIDs(Map(req.Answers, func(answer request.Answer) string { return answer.QuestionId }))
	if err != nil {
		return fmt.Errorf("system cannot find questions for this form: %s", form.Name)
	}

	for _, answer := range req.Answers {
		for _, question := range questions {
			if answer.QuestionId == question.QuestionId.String() {
				var msg *repository.Messaging = nil
				if answer.Messaging != nil {
					msg = &repository.Messaging{
						Email:        answer.Messaging.Email,
						Value3:       answer.Messaging.Value3,
						MessageBox:   answer.Messaging.MessageBox,
						QuestionType: answer.Messaging.QuestionType,
					}
				}
				submissionItems = append(submissionItems, repository.SubmissionDataItem{
					QuestionId: question.QuestionId.String(),
					Question:   question.Question,
					Answer:     answer.Answer,
					Messaging:  msg,
				})
			}
		}
	}

	submissionData := repository.SubmissionData{
		Items: submissionItems,
	}
	createSubmissionParams := repository.CreateSubmissionParams{
		FormId:         form.ID,
		UserId:         req.UserId,
		SubmissionData: submissionData,
		OpenedAt:       req.OpenedAt,
	}
	err = receiver.CreateSubmission(createSubmissionParams)
	if err != nil {
		log.Error("SubmitFormUseCase.answerFormSaveToFormOutputSheet", err)
		return errors.New("system cannot handle the submission")
	}

	defer func() {
		receiver.sendNotification(form)
	}()

	return nil
}

func (receiver *SubmitFormUseCase) sendNotification(form *entity.SForm) {
	questions, err := receiver.GetQuestionsByFormId(form.ID)
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
	if sendNotificationQuestion.QuestionId == "" {
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
		monitor.SendMessageViaTelegram("Failed to send notification for a form submission: ", err.Error())
	}
}
