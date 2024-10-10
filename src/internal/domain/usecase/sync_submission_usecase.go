package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/model"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const timeFormat = "2006-01-02 15:04:05"

var inProcessSubmissionId uint64 = 0

type SyncSubmissionUseCase struct {
	*repository.SubmissionRepository
	*repository.DeviceRepository
	*repository.FormRepository
	*repository.QuestionRepository
	*repository.SettingRepository
	UserSpreadsheetReader    *sheet.Reader
	UserSpreadsheetWriter    *sheet.Writer
	SendEmailUseCase         *SendEmailUseCase
	GetSettingMessageUseCase *GetSettingMessageUseCase
}

func (self *SyncSubmissionUseCase) Execute() {
	if inProcessSubmissionId > 0 {
		monitor.SendMessageViaTelegram("[SKIP][SyncSubmissionUseCase] is running " + strconv.FormatUint(inProcessSubmissionId, 10))
		return
	}
	submission := entity.SSubmission{}
	// Define the GMT+7 timezone
	gmt7 := time.FixedZone("GMT+7", 7*60*60)

	// Get the current time in UTC
	utcNow := time.Now().UTC()

	// Convert the current UTC time to GMT+7
	nowInGMT7 := utcNow.In(gmt7)

	// Parse working hours start and end times
	workingHoursStart, _ := time.Parse("15:04", value.WorkingHoursStart)
	workingHoursEnd, _ := time.Parse("15:04", value.WorkingHoursEnd)

	// Extract the hour and minute from the current time in GMT+7
	currentHour := nowInGMT7.Hour()
	currentMinute := nowInGMT7.Minute()

	// Create a time object for comparison
	currentTime := time.Date(0, 1, 1, currentHour, currentMinute, 0, 0, time.UTC)
	workingStart := time.Date(0, 1, 1, workingHoursStart.Hour(), workingHoursStart.Minute(), 0, 0, time.UTC)
	workingEnd := time.Date(0, 1, 1, workingHoursEnd.Hour(), workingHoursEnd.Minute(), 0, 0, time.UTC)

	// Check if the current time is within working hours
	isWorkingTime := currentTime.After(workingStart) && currentTime.Before(workingEnd)

	if !isWorkingTime {
		log.Info("This time is not working time")
		s, err := self.SubmissionRepository.FindFirstPendingSync()
		//Find first submission

		if err != nil {
			log.Info("Error when find first submission: ", err)
			return
		}

		submission = s
	} else {
		log.Info("This time is working time")
		s, err := self.SubmissionRepository.FindFirstPendingPrioritizedSync()

		if err != nil {
			log.Info("Error when find first prioritized submission: ", err)
			return
		}

		submission = s
	}
	inProcessSubmissionId = submission.ID

	if submission.ID == 0 {
		log.Info("No submission to sync")
		inProcessSubmissionId = 0
		return
	}

	monitor.SendMessageViaTelegram("[INFO] Starting Sync submission (ID): ", strconv.FormatUint(submission.ID, 10),
		fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
		fmt.Sprintf("Submited by %s at %s", submission.DeviceName, submission.CreatedAt.Format("2006-01-02 15:04:05")),
	)

	switch submission.SubmissionType {
	case value.SubmissionTypeValues:
		self.sync(submission)
		break
	case value.SubmissionTypeQrCode:
		self.sync(submission)
		break
	case value.SubmissionTypeTeacher:
		self.sync(submission)
		break
	case value.SubmissionTypeSignUpRegistration:
		self.syncSignUp(submission)
	case value.SubmissionTypeSignUpWriteToMemoryForm:
		self.updateMemoryForm(submission)
	default:
		break
	}
	inProcessSubmissionId = 0
}

func (self *SyncSubmissionUseCase) getFirstOutputRows(spreadSheetID string, sheetName string) ([]string, error) {
	raw, err := self.UserSpreadsheetReader.GetFirstRow(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadSheetID,
		ReadRange:     sheetName + "!I11:WW11",
	})
	if err != nil {
		return nil, err
	}
	rows := make([]string, 0)
	for _, row := range raw {
		if len(row) > 0 {
			rows = append(rows, row[0].(string))
		}
	}

	return rows, err
}

type uniqueQuestionPair struct {
	top    string
	bottom string
}

func (self *SyncSubmissionUseCase) get2FirstOutputRows(spreadSheetID string, sheetName string) ([]uniqueQuestionPair, error) {
	raw, err := self.UserSpreadsheetReader.GetFirstRow(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadSheetID,
		ReadRange:     sheetName + "!I10:WW11",
	})
	if err != nil {
		return nil, err
	}
	questions := make([]uniqueQuestionPair, 0)
	for _, cell := range raw {
		if len(cell) > 1 {
			questions = append(questions, uniqueQuestionPair{
				top:    cell[0].(string),
				bottom: cell[1].(string),
			})
		}
	}

	return questions, err
}

func (self *SyncSubmissionUseCase) checkExistingSignUpForm(question string, existingQuestions []string) bool {
	for _, existingQuestion := range existingQuestions {
		if question == existingQuestion {
			return true
		}
	}
	return false
}

func (self *SyncSubmissionUseCase) checkExisting(questionInForm model.FormQuestionItem, questionsOnSheet []uniqueQuestionPair) bool {
	for _, topRowPair := range questionsOnSheet {
		if questionInForm.QuestionName == topRowPair.bottom {
			//Question Name exists on output form at row 11
			return true
		}
		if questionInForm.QuestionUniqueId != nil && *questionInForm.QuestionUniqueId == topRowPair.bottom {
			//Question Unique ID exists on output form at row 11
			return true
		}
	}
	return false
}

func (self *SyncSubmissionUseCase) findSignUpFormQuestion(question string, questions []model.FormQuestionItem) *model.FormQuestionItem {
	for _, formQuestion := range questions {
		if formQuestion.Question == question {
			return &formQuestion
		}
	}

	return nil
}

func (self *SyncSubmissionUseCase) findFormQuestion(questionPair uniqueQuestionPair, questions []model.FormQuestionItem) *model.FormQuestionItem {
	for _, formQuestion := range questions {
		if formQuestion.Question == questionPair.bottom {
			return &formQuestion
		}
		if formQuestion.QuestionUniqueId != nil && *formQuestion.QuestionUniqueId == questionPair.bottom {
			return &formQuestion
		}
	}

	return nil
}

func (self *SyncSubmissionUseCase) findAnswer(question *model.FormQuestionItem, answers []entity.SubmissionDataItem) *string {
	for _, answer := range answers {
		if answer.QuestionId == question.QuestionId {
			return &answer.Answer
		}
	}
	return nil
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func (self *SyncSubmissionUseCase) sync(submission entity.SSubmission) {
	monitor.SendMessageViaTelegram(
		fmt.Sprintf("[INFO] Sync submission: %d", submission.ID),
		fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
		fmt.Sprintf("Submited by %s at %s", submission.DeviceName, submission.CreatedAt.Format("2006-01-02 15:04:05")),
	)

	var submissionData entity.SubmissionData

	err := json.Unmarshal(submission.SubmissionData, &submissionData)

	if err != nil {
		log.Error("Error when unmarshal submission data: ", err)
		return
	}

	formQuestions, err := self.QuestionRepository.GetQuestionsByFormId(submission.FormId)
	if err != nil {
		log.Error("Error when get questions by form id: ", err)
		return
	}

	//Filter questions by type send_message
	formQuestions = filter(formQuestions, func(q model.FormQuestionItem) bool {
		return value.GetStringValue(value.QuestionSendMessage) != strings.ToLower(q.QuestionType)
	})

	//Filter questions by type send_notification
	formQuestions = filter(formQuestions, func(q model.FormQuestionItem) bool {
		return value.GetStringValue(value.QuestionSendNotification) != strings.ToLower(q.QuestionType)
	})

	firstRows, err := self.get2FirstOutputRows(submission.SpreadsheetId, submission.SheetName)
	if err != nil {
		log.Error("Error when get first rows: ", err, " - ", submission.SpreadsheetId, " - ", submission.SheetName)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][SUBMISSION: %d] Error when get first rows: %s", submission.ID, err.Error()),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("Submited by %s at %s", submission.DeviceName, submission.CreatedAt.Format("2006-01-02 15:04:05")),
			fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
			fmt.Sprintf("Spreadsheet: %s%s", "https://docs.google.com/spreadsheets/d/", submission.SpreadsheetId),
		)
		_ = self.SubmissionRepository.MarkStatusAttempted(submission.ID)
		return
	}

	newQuestions := make([]uniqueQuestionPair, 0)

	for _, formQuestion := range formQuestions {
		if self.checkExisting(formQuestion, firstRows) == false {
			if formQuestion.QuestionUniqueId == nil {
				newQuestions = append(newQuestions, uniqueQuestionPair{
					top:    formQuestion.Question,
					bottom: formQuestion.Question,
				})
			} else {
				newQuestions = append(newQuestions, uniqueQuestionPair{
					top:    formQuestion.Question,
					bottom: *formQuestion.QuestionUniqueId,
				})
			}
		}
	}

	updatingFirst2Rows := make([][]interface{}, 0)
	allQuestions := make([]uniqueQuestionPair, 0)
	for _, row := range firstRows {
		allQuestions = append(allQuestions, row)
	}
	for _, newQuestion := range newQuestions {
		allQuestions = append(allQuestions, newQuestion)
	}
	for _, question := range allQuestions {
		updatingFirst2Rows = append(updatingFirst2Rows, []interface{}{question.top, question.bottom})
	}
	_, err = self.UserSpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
		Range:     submission.SheetName + "!I10",
		Dimension: "COLUMNS",
		Rows:      updatingFirst2Rows,
	}, submission.SpreadsheetId)

	if err != nil {
		log.Error("Error when update range: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][SUBMISSION: %d] Cannot make an update at formQuestion list", submission.ID),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
			fmt.Sprintf("Spreadsheet: %s%s", "https://docs.google.com/spreadsheets/d/", submission.SpreadsheetId),
		)
		return
	}

	//Fix Col from I -> T (Opened At, Submitted At, Timestamp, Device ID, Device Name, Note, INFO1, INFO2, INFO3, FORM Code, FORM Name, FORM Sheet)
	answerRow := make([][]interface{}, 0)
	answerRow = append(answerRow, []interface{}{submission.OpenedAt.Format(timeFormat)})
	answerRow = append(answerRow, []interface{}{submission.CreatedAt.Format(timeFormat)})
	answerRow = append(answerRow, []interface{}{time.Now().Format(timeFormat)})
	answerRow = append(answerRow, []interface{}{submission.DeviceId})
	answerRow = append(answerRow, []interface{}{submission.DeviceName})
	answerRow = append(answerRow, []interface{}{submission.DeviceNote})
	answerRow = append(answerRow, []interface{}{submission.DeviceFirstValue})
	answerRow = append(answerRow, []interface{}{submission.DeviceSecondValue})
	answerRow = append(answerRow, []interface{}{submission.DeviceThirdValue})
	answerRow = append(answerRow, []interface{}{submission.FormNote})
	answerRow = append(answerRow, []interface{}{submission.FormName})
	answerRow = append(answerRow, []interface{}{submission.FormSpreadsheetUrl})

	//var submit []value.Submit = make([]value.Submit, 0)
	//var formSubmission value.FormSubmission = value.FormSubmission{
	//	DeviceName: submission.DeviceName,
	//	Note:       submission.DeviceNote,
	//	Submit:     submit,
	//}

	for index, question := range allQuestions {
		if index > 11 {
			if q := self.findFormQuestion(question, formQuestions); q != nil {
				if ans := self.findAnswer(q, submissionData.Items); ans != nil {
					answerRow = append(answerRow, []interface{}{ans})
				} else {
					answerRow = append(answerRow, []interface{}{""})
				}
			} else {
				answerRow = append(answerRow, []interface{}{""})
			}
		}
	}

	//Get all rows from row number 12
	getAllSubmissionRowsParams := sheet.ReadSpecificRangeParams{
		SpreadsheetId: submission.SpreadsheetId,
		ReadRange:     submission.SheetName + "!I12:WW",
	}
	allSubmissions, err := self.UserSpreadsheetReader.Get(getAllSubmissionRowsParams)
	if err != nil {
		log.Error("Error when get all submission rows: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][SUBMISSION: %d] Cannot get all submission rows", submission.ID),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
		)
		return
	}

	emptyRowIndex := 12
	for _, row := range allSubmissions {
		if len(row) == 0 {
			break
		}
		if len(row) > 2 {
			if row[0] == "" && row[1] == "" && row[2] == "" {
				break
			}
		}
		emptyRowIndex++
	}

	params := sheet.WriteRangeParams{
		Range:     submission.SheetName + "!I" + strconv.Itoa(emptyRowIndex),
		Dimension: "COLUMNS",
		Rows:      answerRow,
	}

	err = self.sendMessage(submission)
	if err != nil {
		log.Error("Error when send message: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR] Failed send message on submission with id: %d", submission.ID),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("Submited by %s at %s", submission.DeviceName, submission.CreatedAt.Format("2006-01-02 15:04:05")),
			fmt.Sprintf("Messaging Data: %v", submission.SubmissionData),
			fmt.Sprintf("Spreadsheet: %s%s", "https://docs.google.com/spreadsheets/d/", submission.SpreadsheetId),
		)
		err = self.SubmissionRepository.MarkStatusSucceeded(submission.ID)
		return
	}

	writtenRanges, err := self.UserSpreadsheetWriter.WriteRangesAsUserEntered(params, submission.SpreadsheetId)
	log.Debug("Wrote: ", writtenRanges)

	if err != nil {
		log.Error("Error when write ranges: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR] Failed to sync submission with id: %d", submission.ID),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("Submited by %s at %s", submission.DeviceName, submission.CreatedAt.Format("2006-01-02 15:04:05")),
			fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
			fmt.Sprintf("Spreadsheet: %s%s", "https://docs.google.com/spreadsheets/d/", submission.SpreadsheetId),
		)
		return
	}

	defer func() {
		self.updateRecentSubmissionToFormSpreadsheet(submission, submissionData)
	}()

	defer func() {
		self.saveCountingDataIfNeeded(submission, submissionData)
	}()

	err = self.SubmissionRepository.MarkStatusSucceeded(submission.ID)
	if err != nil {
		log.Error("Error when mark status succeeded: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR] Failed to mark status succeeded for submission with id: %d", submission.ID),
		)
	} else {
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[SUCCEED] Sync submission with id: %d", submission.ID),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("Submited by %s at %s", submission.DeviceName, submission.CreatedAt.Format("2006-01-02 15:04:05")),
			fmt.Sprintf("Empty row index: %d", emptyRowIndex),
			fmt.Sprintf("Written ranges: %v", writtenRanges),
			fmt.Sprintf("Written At: %s", time.Now().Format("2006-01-02 15:04:05")),
			fmt.Sprintf("Spreadsheet: %s%s", "https://docs.google.com/spreadsheets/d/", submission.SpreadsheetId),
		)
	}
}

func (self *SyncSubmissionUseCase) sendMessage(submission entity.SSubmission) error {
	log.Debug(submission)

	type Messaging struct {
		Email        []string `json:"email" binding:"required"`
		Value3       []string `json:"value3" binding:"required"`
		MessageBox   *string  `json:"messageBox"`
		QuestionType string   `json:"questionType" binding:"required"`
	}

	type ItemElement struct {
		QuestionID string    `json:"question_id" binding:"required"`
		Question   string    `json:"question" binding:"required"`
		Answer     string    `json:"answer" binding:"required"`
		Messaging  Messaging `json:"messaging" binding:"required"`
	}

	type Item struct {
		Items []ItemElement `json:"items" binding:"required"`
	}

	var submissionData Item
	err := json.Unmarshal([]byte(submission.SubmissionData), &submissionData)
	if err != nil {
		return err
	}

	//Filter items by question type send_message
	submissionData.Items = filter(submissionData.Items, func(i ItemElement) bool {
		return i.Messaging.Email != nil
	})

	device, err := self.DeviceRepository.FindDeviceById(submission.DeviceId)
	for _, item := range submissionData.Items {
		bccList := make([]string, 0)
		for _, mailbox := range item.Messaging.Email {
			//validate email using regex
			if _, err := mail.ParseAddress(mailbox); err == nil {
				bccList = append(bccList, mailbox)
			}
		}

		body := ""
		subject := ""

		r, err := self.GetSettingMessageUseCase.Execute(*device)
		if err != nil {
			monitor.SendMessageViaTelegram(
				fmt.Sprintf("[ERROR] Failed to get setting message on submission with id: %d", submission.ID),
			)
			return err
		}
		if len(r.Data.Messages) == 0 {
			monitor.SendMessageViaTelegram(
				fmt.Sprintf("[ERROR] Failed to get setting message on submission with id: %d", submission.ID),
				fmt.Sprintf("to: %v", bccList),
				fmt.Sprintf("error: %s", "No message found"),
			)
			return errors.New("no message found from app setting")
		}
		subject = r.Data.Messages[0].Description
		if item.Messaging.MessageBox != nil {
			body = *item.Messaging.MessageBox + " " + r.Data.Messages[0].Message
		} else {

			body = r.Data.Messages[0].Message
		}

		err = self.SendEmailUseCase.SendMessage(subject, bccList, body)
		if err != nil {
			monitor.SendMessageViaTelegram(
				fmt.Sprintf("[ERROR] Failed to send message on submission with id: %d", submission.ID),
				fmt.Sprintf("to: %v", bccList),
				fmt.Sprintf("error: %s", err.Error()),
			)
		}
	}

	return nil
}

func (self *SyncSubmissionUseCase) syncSignUp(submission entity.SSubmission) {
	monitor.SendMessageViaTelegram(
		fmt.Sprintf("[INFO] Sync submission: %d", submission.ID),
		fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
		fmt.Sprintf("Submited by %s at %s", submission.DeviceName, submission.CreatedAt.Format("2006-01-02 15:04:05")),
	)

	var submissionData entity.SubmissionData

	err := json.Unmarshal([]byte(submission.SubmissionData), &submissionData)

	if err != nil {
		log.Error("Error when unmarshal submission data: ", err)
		return
	}

	formQuestions, err := self.QuestionRepository.GetQuestionsByFormId(submission.FormId)
	if err != nil {
		log.Error("Error when get questions by form id: ", err)
		return
	}

	firstRows, err := self.getFirstOutputRows(submission.SpreadsheetId, submission.SheetName)
	if err != nil {
		log.Error("Error when get first rows: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][SUBMISSION: %d] Error when get first rows: %s", submission.ID, err.Error()),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("Submited by %s at %s", submission.DeviceName, submission.CreatedAt.Format("2006-01-02 15:04:05")),
			fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
		)
		_ = self.SubmissionRepository.MarkStatusAttempted(submission.ID)
		return
	}

	newQuestions := make([]interface{}, 0)
	for _, question := range formQuestions {
		if self.checkExistingSignUpForm(question.Question, firstRows) == false {
			newQuestions = append(newQuestions, question.Question)
		}
	}

	newFirstRows := make([][]interface{}, 0)
	allQuestions := make([]string, 0)
	for _, row := range firstRows {
		newFirstRows = append(newFirstRows, []interface{}{row})
		allQuestions = append(allQuestions, row)
	}
	for _, newQuestion := range newQuestions {
		newFirstRows = append(newFirstRows, []interface{}{newQuestion})
		allQuestions = append(allQuestions, newQuestion.(string))
	}
	_, err = self.UserSpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
		Range:     submission.SheetName + "!I11",
		Dimension: "COLUMNS",
		Rows:      newFirstRows,
	}, submission.SpreadsheetId)

	if err != nil {
		log.Error("Error when update range: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][SUBMISSION: %d] Cannot make an update at question list", submission.ID),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
		)
		return
	}

	answerRow := make([][]interface{}, 0)
	answerRow = append(answerRow, []interface{}{submission.OpenedAt.Format(timeFormat)})
	answerRow = append(answerRow, []interface{}{submission.CreatedAt.Format(timeFormat)})
	answerRow = append(answerRow, []interface{}{time.Now().Format(timeFormat)})
	answerRow = append(answerRow, []interface{}{submission.DeviceId})
	answerRow = append(answerRow, []interface{}{submission.DeviceName})
	answerRow = append(answerRow, []interface{}{submission.DeviceNote})
	answerRow = append(answerRow, []interface{}{submission.DeviceFirstValue})
	answerRow = append(answerRow, []interface{}{submission.DeviceSecondValue})
	answerRow = append(answerRow, []interface{}{submission.DeviceThirdValue})
	answerRow = append(answerRow, []interface{}{submission.FormNote})
	answerRow = append(answerRow, []interface{}{submission.FormName})
	answerRow = append(answerRow, []interface{}{submission.FormSpreadsheetUrl})

	//var submit []value.Submit = make([]value.Submit, 0)
	//var formSubmission value.FormSubmission = value.FormSubmission{
	//	DeviceName: submission.DeviceName,
	//	Note:       submission.DeviceNote,
	//	Submit:     submit,
	//}

	for index, question := range allQuestions {
		if index > 11 {
			if q := self.findSignUpFormQuestion(question, formQuestions); q != nil {
				answerRow = append(answerRow, []interface{}{self.findAnswer(q, submissionData.Items)})
				//formSubmission.Submit = append(formSubmission.Submit, value.Submit{
				//	Question: q.Question,
				//	Answer:   self.findAnswer(q.QuestionId, submissionData.Items),
				//})
			} else {
				answerRow = append(answerRow, []interface{}{""})
			}
		}
	}

	//Get all rows from row number 12
	getAllSubmissionRowsParams := sheet.ReadSpecificRangeParams{
		SpreadsheetId: submission.SpreadsheetId,
		ReadRange:     submission.SheetName + "!I12:WW",
	}
	allSubmissions, err := self.UserSpreadsheetReader.Get(getAllSubmissionRowsParams)
	if err != nil {
		log.Error("Error when get all submission rows: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][SUBMISSION: %d] Cannot get all submission rows", submission.ID),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
		)
		return
	}

	emptyRowIndex := 12
	for _, row := range allSubmissions {
		if len(row) == 0 {
			break
		}
		if len(row) > 2 {
			if row[0] == "" && row[1] == "" && row[2] == "" {
				break
			}
		}
		emptyRowIndex++
	}

	params := sheet.WriteRangeParams{
		Range:     submission.SheetName + "!I" + strconv.Itoa(emptyRowIndex),
		Dimension: "COLUMNS",
		Rows:      answerRow,
	}

	writtenRanges, err := self.UserSpreadsheetWriter.WriteRangesAsUserEntered(params, submission.SpreadsheetId)
	log.Debug("Wrote :", writtenRanges)

	if err != nil {
		log.Error("Error when write ranges: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR] Failed to sync submission with id: %d", submission.ID),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("Submited by %s at %s", submission.DeviceName, submission.CreatedAt.Format("2006-01-02 15:04:05")),
			fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
		)
		return
	}

	err = self.SubmissionRepository.MarkStatusSucceeded(submission.ID)
	if err != nil {
		log.Error("Error when mark status succeeded: ", err)
	} else {
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[SUCCEED] Sync submission with id: %d", submission.ID),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("Submited by %s at %s", submission.DeviceName, submission.CreatedAt.Format("2006-01-02 15:04:05")),
			fmt.Sprintf("Empty row index: %d", emptyRowIndex),
			fmt.Sprintf("Written ranges: %v", writtenRanges),
			fmt.Sprintf("Written At: %s", time.Now().Format("2006-01-02 15:04:05")),
		)
	}
}

func (self *SyncSubmissionUseCase) checkExisting2(question string, existingQuestions []string) bool {
	for _, existingQuestion := range existingQuestions {
		if question == existingQuestion {
			return true
		}
	}
	return false
}

func (self *SyncSubmissionUseCase) updateRecentSubmissionToFormSpreadsheet(submission entity.SSubmission, submissionData entity.SubmissionData) {
	form, err := self.FormRepository.GetFormById(submission.FormId)
	if err != nil {
		log.Error("Form ", submission.FormId, " does not exist.")
		return
	}

	if form.Type != value.FormType_SelfRemember {
		log.Debug("No need update recent submission for the this form")
		return
	}

	formQuestionKeys, err := self.UserSpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: form.SpreadsheetId,
		ReadRange:     form.SheetName + "!K13:L",
	})

	if err != nil {
		log.Error("Could not read form ", form.SpreadsheetId, " at sheet ", form.SheetName, " raw error ", err.Error())
		return
	}

	recentSubmit := make([][]interface{}, 0)

	for i := 0; i < len(formQuestionKeys); i++ {
		if len(formQuestionKeys[i]) != 0 && formQuestionKeys[i][0].(string) != "" {
			answer := ""
			for j := 0; j < len(submissionData.Items); j++ {
				if strings.ToLower(submissionData.Items[j].QuestionId) == strings.ToLower((form.Note + "_" + form.SpreadsheetId + "_" + formQuestionKeys[i][0].(string))) {
					answer = submissionData.Items[j].Answer
				}
			}
			recentSubmit = append(recentSubmit, []interface{}{answer})
		} else {
			recentSubmit = append(recentSubmit, []interface{}{""})
		}
	}

	params := sheet.WriteRangeParams{
		Range:     form.SheetName + "!H13:H",
		Dimension: "ROWS",
		Rows:      recentSubmit,
	}

	_, err = self.UserSpreadsheetWriter.UpdateRange(params, form.SpreadsheetId)

	if err != nil {
		log.Error(err.Error())
		monitor.SendMessageViaTelegram("Update recent submission for memory form failed", err.Error())
	}
}

func (self *SyncSubmissionUseCase) findFormQuestion2(question string, questions []model.FormQuestionItem) *model.FormQuestionItem {
	for _, formQuestion := range questions {
		if formQuestion.Question == question {
			return &formQuestion
		}
	}
	return nil
}

func (self *SyncSubmissionUseCase) updateMemoryForm(submission entity.SSubmission) {
	monitor.SendMessageViaTelegram(
		fmt.Sprintf("[INFO] Sync submission: %d", submission.ID),
		fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
		fmt.Sprintf("Submited by %s at %s", submission.DeviceName, submission.CreatedAt.Format("2006-01-02 15:04:05")),
	)

	var submissionData entity.SubmissionData

	err := json.Unmarshal([]byte(submission.SubmissionData), &submissionData)

	if err != nil {
		log.Error("Error when unmarshal submission data: ", err)
		return
	}

	formQuestions, err := self.QuestionRepository.GetQuestionsByFormId(submission.FormId)
	if err != nil {
		log.Error("Error when get questions by form id: ", err)
		return
	}

	firstRows, err := self.getFirstOutputRows(submission.SpreadsheetId, submission.SheetName)
	if err != nil {
		log.Error("Error when get first rows: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][SUBMISSION: %d] Error when get first rows: %s", submission.ID, err.Error()),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("Submited by %s at %s", submission.DeviceName, submission.CreatedAt.Format("2006-01-02 15:04:05")),
			fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
		)
		_ = self.SubmissionRepository.MarkStatusAttempted(submission.ID)
		return
	}

	newQuestions := make([]interface{}, 0)
	for _, question := range formQuestions {
		if self.checkExisting2(question.Question, firstRows) == false {
			newQuestions = append(newQuestions, question.Question)
		}
	}

	newFirstRows := make([][]interface{}, 0)
	allQuestions := make([]string, 0)
	for _, row := range firstRows {
		newFirstRows = append(newFirstRows, []interface{}{row})
		allQuestions = append(allQuestions, row)
	}
	for _, newQuestion := range newQuestions {
		newFirstRows = append(newFirstRows, []interface{}{newQuestion})
		allQuestions = append(allQuestions, newQuestion.(string))
	}
	_, err = self.UserSpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
		Range:     submission.SheetName + "!I11",
		Dimension: "COLUMNS",
		Rows:      newFirstRows,
	}, submission.SpreadsheetId)

	if err != nil {
		log.Error("Error when update range: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][SUBMISSION: %d] Cannot make an update at question list", submission.ID),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
		)
		return
	}

	answerRow := make([][]interface{}, 0)
	answerRow = append(answerRow, []interface{}{submission.OpenedAt.Format(timeFormat)})
	answerRow = append(answerRow, []interface{}{submission.CreatedAt.Format(timeFormat)})
	answerRow = append(answerRow, []interface{}{time.Now().Format(timeFormat)})
	answerRow = append(answerRow, []interface{}{submission.DeviceId})
	answerRow = append(answerRow, []interface{}{submission.DeviceName})
	answerRow = append(answerRow, []interface{}{submission.DeviceNote})
	answerRow = append(answerRow, []interface{}{submission.DeviceFirstValue})
	answerRow = append(answerRow, []interface{}{submission.DeviceSecondValue})
	answerRow = append(answerRow, []interface{}{submission.DeviceThirdValue})
	answerRow = append(answerRow, []interface{}{submission.FormNote})
	answerRow = append(answerRow, []interface{}{submission.FormName})
	answerRow = append(answerRow, []interface{}{submission.FormSpreadsheetUrl})

	for index, question := range allQuestions {
		if index > 11 {
			if q := self.findFormQuestion2(question, formQuestions); q != nil {
				answerRow = append(answerRow, []interface{}{self.findAnswer(q, submissionData.Items)})
				//formSubmission.Submit = append(formSubmission.Submit, value.Submit{
				//	Question: q.Question,
				//	Answer:   self.findAnswer(q.QuestionId, submissionData.Items),
				//})
			} else {
				answerRow = append(answerRow, []interface{}{""})
			}
		}
	}

	//Get all rows from row number 12
	getAllSubmissionRowsParams := sheet.ReadSpecificRangeParams{
		SpreadsheetId: submission.SpreadsheetId,
		ReadRange:     submission.SheetName + "!I12:WW",
	}
	allSubmissions, err := self.UserSpreadsheetReader.Get(getAllSubmissionRowsParams)
	if err != nil {
		log.Error("Error when get all submission rows: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][SUBMISSION: %d] Cannot get all submission rows", submission.ID),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
		)
		return
	}

	emptyRowIndex := 12
	for _, row := range allSubmissions {
		if len(row) == 0 {
			break
		}
		if len(row) > 2 {
			if row[0] == "" && row[1] == "" && row[2] == "" {
				break
			}
		}
		emptyRowIndex++
	}

	self.updateRecentSubmissionToFormSpreadsheet(submission, submissionData)

	err = self.MarkStatusSucceeded(submission.ID)
	if err != nil {
		log.Error("Error when mark status succeeded: ", err)
	}
}
