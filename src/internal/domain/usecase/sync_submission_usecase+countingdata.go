package usecase

import (
	"encoding/json"
	"regexp"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/model"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/sheet"
	"time"

	log "github.com/sirupsen/logrus"
)

// Write new record in to device uploader at CountingData sheet
func (receiver *SyncSubmissionUseCase) saveCountingDataIfNeeded(submission entity.SSubmission, submissionData entity.SubmissionData) {
	var codeCountingQuestion model.FormQuestionItem
	questions, err := receiver.QuestionRepository.GetQuestionsByFormId(submission.FormId)
	if err != nil {
		log.Error("Error when get questions: ", err)
		return
	}
	for _, q := range questions {
		if q.QuestionType == value.GetRawValue(value.QuestionCodeCounting) {
			codeCountingQuestion = q
		}
	}

	if codeCountingQuestion.QuestionId == "" {
		log.Debug("No Code Generator question in the form question")
		return
	}

	log.Debug("codeCountingQuestion: ", codeCountingQuestion)

	// Find the code value from the submission
	var codeValue string
	for _, item := range submissionData.Items {
		if item.QuestionId == codeCountingQuestion.QuestionId {
			codeValue = item.Answer
			break
		}
	}

	if codeValue == "" {
		log.Warning("No code value in the submission")
		return
	}

	//Extract Prefix from Attributes
	type Attr struct {
		Value string `json:"value"`
	}
	var questionAttributes = Attr{}
	err = json.Unmarshal(codeCountingQuestion.Attributes, &questionAttributes)

	if err != nil {
		log.Error("failed to unmarshal code counting question attributes")
		return
	}

	if questionAttributes.Value == "" {
		log.Warning("No code value in the question attributes")
		return
	}

	codeCountingSettingData, err := receiver.SettingRepository.GetCodeCountingDataSetting()
	if err != nil {
		log.Error("failed to get sync devices settings")
		return
	}

	type APIDistributerSetting struct {
		SettingName    string `json:"setting_name"`
		SpreadSheetUrl string `json:"spreadsheet_url"`
	}

	var codeCountingSettings APIDistributerSetting

	err = json.Unmarshal(codeCountingSettingData.Settings, &codeCountingSettings)
	if err != nil {
		log.Info(err.Error())
	}
	codeCountingSettings.SettingName = codeCountingSettingData.SettingName

	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(codeCountingSettings.SpreadSheetUrl)

	if len(match) < 2 {
		log.Error("failed to parse spreadsheet id from sync devices sheet")
		return
	}

	spreadsheetId := match[1]

	//Get data at row 11 from the spreadsheet
	row11Data, err := receiver.UserSpreadsheetReader.GetFirstRow(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "CountingData!U11:Z11",
	})

	if err != nil {
		log.Error("Error when get row 11 data from counting data sheet: ", err)
		// monitor.SendMessageViaTelegram(
		// 	fmt.Sprintf("[ERROR][SUBMISSION: %d] Cannot get row 11 data", submission.ID),
		// 	fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
		// 	fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
		// )
		return
	}

	timeFormat := "2006-01-02T15:04:05Z"

	////Append device's values
	//rows = append(rows, []interface{}{submission.DeviceFirstValue})
	//rows = append(rows, []interface{}{submission.DeviceSecondValue})
	//rows = append(rows, []interface{}{submission.DeviceThirdValue})
	//
	////Append Form's values
	//rows = append(rows, []interface{}{submission.FormName})
	//rows = append(rows, []interface{}{submission.FormSpreadsheetUrl})
	//rows = append(rows, []interface{}{submission.CreatedAt.Format(timeFormat)})
	//
	////Append Question
	//rows = append(rows, []interface{}{codeValue})

	//_, err = receiver.UserSpreadsheetWriter.WriteRanges(sheet.WriteRangeParams{
	//	Range:     "CountingData!K11",
	//	Rows:      rows,
	//	Dimension: "COLUMNS",
	//}, spreadsheetId)
	//if err != nil {
	//	log.Error("Error when append counting data: ", err)
	//	monitor.SendMessageViaTelegram(
	//		fmt.Sprintf("[ERROR][SUBMISSION: %d] Cannot append new record to counting data", submission.ID),
	//		fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
	//		fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
	//	)
	//	return
	//}

	answerRow := make([][]interface{}, 0)
	answerRow = append(answerRow, []interface{}{submission.OpenedAt.Format(timeFormat)})
	answerRow = append(answerRow, []interface{}{submission.CreatedAt.Format(timeFormat)})
	answerRow = append(answerRow, []interface{}{time.Now().Format(timeFormat)})
	answerRow = append(answerRow, []interface{}{submission.DeviceId})
	answerRow = append(answerRow, []interface{}{nil})
	answerRow = append(answerRow, []interface{}{nil})
	answerRow = append(answerRow, []interface{}{nil})
	answerRow = append(answerRow, []interface{}{nil})
	answerRow = append(answerRow, []interface{}{nil})
	answerRow = append(answerRow, []interface{}{nil})
	answerRow = append(answerRow, []interface{}{nil})
	answerRow = append(answerRow, []interface{}{nil})

	//Find the code from row 11 begin at column U
	isExist := false
	newRow11Data := make([][]interface{}, 0)
	if row11Data == nil {
		newRow11Data = append(newRow11Data, []interface{}{questionAttributes.Value})
		answerRow = append(answerRow, []interface{}{codeValue})
	} else {
		for index, row := range row11Data {
			if len(row) == 0 {
				newRow11Data = append(newRow11Data, []interface{}{questionAttributes.Value})
				answerRow = append(answerRow, []interface{}{codeValue})
				break
			} else if len(row) > 0 && row[0] == questionAttributes.Value {
				isExist = true
				answerRow = append(answerRow, []interface{}{codeValue})
				break
			} else if index == len(row11Data)-1 {
				newRow11Data = append(newRow11Data, row)
				answerRow = append(answerRow, []interface{}{nil})
				newRow11Data = append(newRow11Data, []interface{}{questionAttributes.Value})
				answerRow = append(answerRow, []interface{}{codeValue})
				break
			} else {
				newRow11Data = append(newRow11Data, row)
				answerRow = append(answerRow, []interface{}{nil})
			}
		}
	}

	if !isExist {
		//Update row 11 a column U
		_, err = receiver.UserSpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
			Range:     "CountingData!U11",
			Dimension: "COLUMNS",
			Rows:      newRow11Data,
		}, spreadsheetId)
		if err != nil {
			log.Error("Error when update row 11 data: ", err)
			// monitor.SendMessageViaTelegram(
			// 	fmt.Sprintf("[ERROR][SUBMISSION: %d] Cannot update row 11 data", submission.ID),
			// 	fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
			// 	fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
			// )
			return
		}
	}

	_, err = receiver.UserSpreadsheetWriter.WriteRanges(sheet.WriteRangeParams{
		Range:     "CountingData!K11",
		Rows:      answerRow,
		Dimension: "COLUMNS",
	}, spreadsheetId)
	if err != nil {
		log.Error("Error when append counting data: ", err)
		// monitor.SendMessageViaTelegram(
		// 	fmt.Sprintf("[ERROR][SUBMISSION: %d] Cannot append new record to counting data", submission.ID),
		// 	fmt.Sprintf("Target Form Note %s - Form Name %s", submission.FormNote, submission.FormName),
		// 	fmt.Sprintf("[GOOGLE API]: %s", err.Error()),
		// )
		return
	}
}
