package usecase

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/model"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"
	"time"
)

func (self *SyncSubmissionUseCase) checkExistingUniqueId(questionUniqueID *string, existingQuestions []uniqueQuestionPair) bool {
	if questionUniqueID == nil {
		return false
	}
	for _, existingQuestion := range existingQuestions {
		if *questionUniqueID == existingQuestion.bottom {
			return true
		}
	}
	return false
}

func (self *SyncSubmissionUseCase) syncSubmissionContainsQuestionUniqueId(submission entity.SSubmission, submissionData entity.SubmissionData, formQuestions []model.FormQuestionItem) {
	//Filter questions by type send_message
	formQuestions = filter(formQuestions, func(q model.FormQuestionItem) bool {
		return value.GetStringValue(value.QuestionSendMessage) != strings.ToLower(q.QuestionType)
	})

	firstRows, err := self.get2FirstOutputRows(submission.SpreadsheetId, submission.SheetName)
	if err != nil {
		log.Error("Error when get first rows: ", err)
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
	for _, question := range formQuestions {
		if self.checkExistingUniqueId(question.QuestionUniqueId, firstRows) == false {
			newQuestions = append(newQuestions, uniqueQuestionPair{
				top:    question.Question,
				bottom: *question.QuestionUniqueId,
			})
		}
	}

	newFirstRows := make([][]interface{}, 0)
	allQuestions := make([]uniqueQuestionPair, 0)
	for _, row := range firstRows {
		allQuestions = append(allQuestions, row)
	}
	for _, newQuestion := range newQuestions {
		allQuestions = append(allQuestions, newQuestion)
	}

	for _, question := range allQuestions {
		newFirstRows = append(newFirstRows, []interface{}{question.top, question.bottom})
	}

	_, err = self.UserSpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
		Range:     submission.SheetName + "!I10",
		Dimension: "COLUMNS",
		Rows:      newFirstRows,
	}, submission.SpreadsheetId)

	if err != nil {
		log.Error("Error when update range: ", err)
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][SUBMISSION: %d] Cannot make an update at question list", submission.ID),
			fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
			fmt.Sprintf("Spreadsheet: %s%s", "https://docs.google.com/spreadsheets/d/", submission.SpreadsheetId),
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
			if q := self.findFormQuestionByUniqueId(question.bottom, formQuestions); q != nil {
				answerRow = append(answerRow, []interface{}{
					self.findAnswerForUniqueIDQuestion(q, submissionData.Items),
				})
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

func (self *SyncSubmissionUseCase) findFormQuestionByUniqueId(UniqueId string, questions []model.FormQuestionItem) *model.FormQuestionItem {
	for _, formQuestion := range questions {
		if formQuestion.QuestionUniqueId != nil && *formQuestion.QuestionUniqueId == UniqueId {
			return &formQuestion
		}
	}
	return nil
}

func (self *SyncSubmissionUseCase) findAnswerForUniqueIDQuestion(question *model.FormQuestionItem, answers []entity.SubmissionDataItem) string {
	for _, answer := range answers {
		if answer.QuestionId == question.QuestionId {
			return answer.Answer
		}
	}
	return ""
}
