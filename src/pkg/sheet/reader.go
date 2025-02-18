package sheet

import (
	"errors"
	"strconv"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

type Reader struct {
	sheetsService *sheets.Service
}

// / ReadSpecificRangeParams is the params for reading a specific range of a spreadsheet
// / SpreadsheetId is the id of the spreadsheet
// / ReadRange is the range to read. Eg, "Sheet1!A1:B2"
type ReadSpecificRangeParams struct {
	SpreadsheetId string
	ReadRange     string
}

// / ReadSpecificRange reads a specific range of a spreadsheet
// / params is the params for reading a specific range of a spreadsheet
// / returns the rows of the spreadsheet
func (receiver Reader) Get(params ReadSpecificRangeParams) ([][]interface{}, error) {
	//monitor.SendMessageViaTelegram("[GOOGLE API]Reading sheet " + params.SpreadsheetId + " " + params.ReadRange)
	resp, err := receiver.sheetsService.Spreadsheets.Values.Get(params.SpreadsheetId, params.ReadRange).
		ValueRenderOption("FORMATTED_VALUE").
		Do()
	if err != nil {
		log.Error("Unable to retrieve data from sheet:", err)
		return nil, err
	}

	if len(resp.Values) == 0 {
		log.Info("No data found.")
		return nil, err
	} else {
		return resp.Values, nil
	}
}

// / ReadSpecificRange reads a specific range of a spreadsheet
// / params is the params for reading a specific range of a spreadsheet
// / returns the rows of the spreadsheet
func (receiver Reader) GetFirstRow(params ReadSpecificRangeParams) ([][]interface{}, error) {
	//monitor.SendMessageViaTelegram("[GOOGLE API]GetFirstRow " + params.SpreadsheetId + " " + params.ReadRange)
	resp, err := receiver.sheetsService.Spreadsheets.Values.Get(params.SpreadsheetId, params.ReadRange).
		MajorDimension("COLUMNS").
		ValueRenderOption("FORMATTED_VALUE").
		Do()
	if err != nil {
		log.Error("Unable to retrieve data from sheet:", err)
		return nil, err
	}

	if len(resp.Values) == 0 {
		log.Info("No data found.")
		return nil, err
	} else {
		return resp.Values, nil
	}
}

func (receiver Reader) FindFirstRow(params ReadSpecificRangeParams, deviceID string) (int, error) {
	//monitor.SendMessageViaTelegram("[GOOGLE API]FindFirstRow " + params.SpreadsheetId + " " + params.ReadRange)
	resp, err := receiver.sheetsService.Spreadsheets.Values.Update(params.SpreadsheetId, params.ReadRange, &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         [][]interface{}{{"=MATCH(\"" + deviceID + "\", Devices!L:L, 0)"}},
	}).ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		log.Error("Unable to retrieve data from sheet:", err)
		return 0, err
	}

	updatedRows, err := receiver.sheetsService.Spreadsheets.Values.Get(params.SpreadsheetId, "LOOKUP_SHEET!A1").MajorDimension("COLUMNS").
		ValueRenderOption("FORMATTED_VALUE").Do()
	if err != nil {
		log.Error("Unable to retrieve data from sheet:", err)
		return 0, err
	}

	rowNo := 0
	if len(updatedRows.Values) > 0 {
		if len(updatedRows.Values[0]) > 0 {
			rowNoInString := updatedRows.Values[0][0].(string)
			rowNo, err = strconv.Atoi(rowNoInString)
		}
	}
	if rowNo == 0 {
		return 0, errors.New("Unable to find row number for device " + deviceID)
	}

	log.Debug("Wrote: ", resp)
	return rowNo, err
}

func (receiver Reader) GetSheets(spreadsheetId string) ([]string, error) {
	//monitor.SendMessageViaTelegram("[GOOGLE API]GetSheets " + spreadsheetId)
	resp, err := receiver.sheetsService.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		log.Error("Unable to retrieve data from sheet:", err)
		return nil, err
	}

	sheetNames := make([]string, len(resp.Sheets))
	for i, sheet := range resp.Sheets {
		sheetNames[i] = sheet.Properties.Title
	}

	return sheetNames, nil
}

type SingleSheet struct {
	ID    int64
	Title string
}

func (receiver Reader) GetAllSheets(spreadsheetId string) ([]SingleSheet, error) {
	//monitor.SendMessageViaTelegram("[GOOGLE API]GetSheets " + spreadsheetId)
	resp, err := receiver.sheetsService.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		log.Error("Unable to retrieve data from sheet:", err)
		return nil, err
	}

	singleSheets := make([]SingleSheet, 0)
	for _, sheet := range resp.Sheets {
		if sheet.Properties != nil {
			singleSheets = append(singleSheets, SingleSheet{
				ID:    sheet.Properties.SheetId,
				Title: sheet.Properties.Title,
			})
		}
	}

	return singleSheets, nil
}
