package sheet

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

type Writer struct {
	sheetsService *sheets.Service
}

type WriteRangeParams struct {
	Range     string
	Dimension string
	Rows      [][]interface{}
}

type ClearRangeParams struct {
	SpreadsheetId string
	Range         string
}

type AppendParams struct {
	SheetName string
	Dimension string
	Rows      [][]interface{}
}

func (receiver Writer) WriteRanges(params WriteRangeParams, spreadsheetId string) (*sheets.AppendValuesResponse, error) {
	//monitor.SendMessageViaTelegram("[GOOGLE API]Writing sheet " + spreadsheetId + " - Append at range " + params.Range)
	var updateValues = &sheets.ValueRange{
		MajorDimension: params.Dimension,
		Range:          params.Range,
		Values:         params.Rows,
	}
	resp, err := receiver.sheetsService.Spreadsheets.Values.Append(spreadsheetId, params.Range, updateValues).ValueInputOption("RAW").Do()
	if err != nil {
		log.Error("Unable to append data from sheet: ", err)
		return nil, err
	}

	log.Debug("Wrote: ", resp)

	return resp, nil
}

func (receiver Writer) WriteRangesAsUserEntered(params WriteRangeParams, spreadsheetId string) (*sheets.AppendValuesResponse, error) {
	//monitor.SendMessageViaTelegram("[GOOGLE API]Writing sheet " + spreadsheetId + " - Append at range " + params.Range)
	var updateValues = &sheets.ValueRange{
		MajorDimension: params.Dimension,
		Range:          params.Range,
		Values:         params.Rows,
	}
	resp, err := receiver.sheetsService.Spreadsheets.Values.Append(spreadsheetId, params.Range, updateValues).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Error("Unable to append data from sheet: ", err)
		return nil, err
	}

	log.Debug("Wrote: ", resp)

	return resp, nil
}

func (receiver Writer) UpdateRange(params WriteRangeParams, spreadsheetId string) (*sheets.UpdateValuesResponse, error) {
	//monitor.SendMessageViaTelegram("[GOOGLE API]Writing sheet " + spreadsheetId + " - Update at range " + params.Range)
	var updateValues = &sheets.ValueRange{
		MajorDimension: params.Dimension,
		Range:          params.Range,
		Values:         params.Rows,
	}
	resp, err := receiver.sheetsService.Spreadsheets.Values.Update(spreadsheetId, params.Range, updateValues).ValueInputOption("RAW").Do()
	if err != nil {
		log.Error("Unable to retrieve data from sheet: ", err)
		return nil, err
	}

	log.Debug("Wrote: ", resp)

	return resp, nil
}

func (receiver Writer) UpdateRanges(spreadsheetId string, params []WriteRangeParams) error {
	//monitor.SendMessageViaTelegram("[GOOGLE API]Writing sheet " + spreadsheetId + " - Update at ranges ")
	rbb := &sheets.BatchUpdateSpreadsheetRequest{}
	//for _, p := range params {
	//	rbb.Requests = append(rbb.Requests, &sheets.Request{
	//		UpdateCells: &sheets.UpdateCellsRequest{
	//			Fields: "*",
	//			Rows: []*sheets.RowData{
	//				{
	//					Values: []*sheets.CellData{
	//						{
	//							UserEnteredValue: &sheets.ExtendedValue{
	//								StringValue: "test",
	//							},
	//						},
	//					},
	//				},
	//			},
	//			Start: &sheets.GridCoordinate{
	//				RowIndex:    0,
	//				ColumnIndex: 0,
	//			},
	//		},
	//	})
	//}

	resp, err := receiver.sheetsService.Spreadsheets.BatchUpdate(spreadsheetId, rbb).Do()
	if err != nil {
		return err
	}

	log.Debug("UpdateRanges: ", resp)

	return err
}

func (receiver Writer) CreateSheet(sheetName string, spreadsheetId string) error {
	//monitor.SendMessageViaTelegram("[GOOGLE API]Creating sheet " + spreadsheetId + " - sheet name " + sheetName)
	req := sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: sheetName,
			},
		},
	}

	rbb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{&req},
	}

	resp, err := receiver.sheetsService.Spreadsheets.BatchUpdate(spreadsheetId, rbb).Do()
	if err != nil {
		return err
	}

	log.Debug("Create new Sheet", resp)

	return nil
}

func (receiver Writer) AppendSheet(params AppendParams, spreadsheetId string) (*sheets.UpdateValuesResponse, error) {
	//monitor.SendMessageViaTelegram("[GOOGLE API]Writing sheet " + spreadsheetId + " - Append at sheet " + params.SheetName)
	log.Debug("AppendSheet: ", params)
	log.Debug("AppendSheet: ", spreadsheetId)
	data := &sheets.ValueRange{
		Range:          params.SheetName,
		MajorDimension: params.Dimension,
		Values:         params.Rows,
	}

	resp, err := receiver.sheetsService.Spreadsheets.Values.Append(spreadsheetId, params.SheetName, data).ValueInputOption("RAW").Do()
	if err != nil {
		return nil, err
	}

	log.Debug("Append Sheet", resp)

	return resp.Updates, nil
}

func (receiver Writer) ClearRange(params ClearRangeParams) (*sheets.ClearValuesResponse, error) {
	log.Debug("ClearRange: ", params)
	resp, err := receiver.sheetsService.Spreadsheets.Values.
		Clear(params.SpreadsheetId, params.Range, &sheets.ClearValuesRequest{}).
		Do()

	if err != nil {
		return nil, err
	}

	log.Debug("Clear Ranges: ", resp)

	return resp, nil
}

type CopySingleSheetParam struct {
	FromSpreadsheetId string
	SingleSheet       SingleSheet
	ToSpreadsheetId   string
}

func (receiver Writer) CopySingleSheet(params CopySingleSheetParam) error {
	copyRequest := &sheets.CopySheetToAnotherSpreadsheetRequest{
		DestinationSpreadsheetId: params.ToSpreadsheetId,
	}
	copiedSheet, err := receiver.sheetsService.Spreadsheets.Sheets.
		CopyTo(params.FromSpreadsheetId, params.SingleSheet.ID, copyRequest).
		Context(context.Background()).
		Do()

	if err != nil {
		return err
	}

	log.Debug(copiedSheet)

	request := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
					Properties: &sheets.SheetProperties{
						SheetId: copiedSheet.SheetId,
						Title:   params.SingleSheet.Title,
					},
					Fields: "title",
				},
			},
		},
	}

	_, err = receiver.sheetsService.Spreadsheets.
		BatchUpdate(params.ToSpreadsheetId, request).
		Context(context.Background()).
		Do()
	if err != nil {
		log.Fatalf("Unable to rename sheet: %v", err)
	}

	return err
}

type DeleteSheetParams struct {
	SpreadsheetID string
	SheetTitle    string
}

func (receiver Writer) DeleteSheet(params DeleteSheetParams) error {
	resp, err := receiver.sheetsService.Spreadsheets.Get(params.SpreadsheetID).Do()
	if err != nil {
		log.Error("Unable to retrieve data from sheet:", err)
		return err
	}

	for _, sheet := range resp.Sheets {
		if sheet.Properties != nil && sheet.Properties.Title == params.SheetTitle {
			// Create a batch update request to delete the sheet
			request := &sheets.BatchUpdateSpreadsheetRequest{
				Requests: []*sheets.Request{
					{
						DeleteSheet: &sheets.DeleteSheetRequest{
							SheetId: sheet.Properties.SheetId,
						},
					},
				},
			}

			_, err = receiver.sheetsService.Spreadsheets.
				BatchUpdate(params.SpreadsheetID, request).
				Context(context.Background()).
				Do()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type DuplicateSpreadsheetParams struct {
	SourceSpreadsheetId   string
	TargetSpreadsheetName string
	TargetSheetName       string
}

type DuplicateSpreadsheetResult struct {
	SpreadsheetId string
}

func (receiver Writer) DuplicateSpreadsheet(params DuplicateSpreadsheetParams) (DuplicateSpreadsheetResult, error) {
	ctx := context.Background()
	resp, err := receiver.sheetsService.Spreadsheets.
		Get(params.SourceSpreadsheetId).
		Context(ctx).
		Do()
	if err != nil {
		log.WithError(err).Error("Failed to retrieve spreadsheet")
		return DuplicateSpreadsheetResult{}, err
	}

	var signUpSheetId int64

	for _, sheet := range resp.Sheets {
		if sheet.Properties.Title == params.TargetSheetName {
			signUpSheetId = sheet.Properties.SheetId
		}
	}

	if signUpSheetId == 0 {
		log.WithError(err).Error("Failed to find sheet")
		return DuplicateSpreadsheetResult{}, err
	}

	createResp, err := receiver.sheetsService.Spreadsheets.Create(&sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: params.TargetSpreadsheetName,
		},
	}).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Failed to create new spreadsheet: %v", err)
	}

	destinationSpreadsheetId := createResp.SpreadsheetId

	err = receiver.CopySingleSheet(CopySingleSheetParam{
		FromSpreadsheetId: params.SourceSpreadsheetId,
		SingleSheet: SingleSheet{
			ID:    signUpSheetId,
			Title: params.TargetSheetName,
		},
		ToSpreadsheetId: destinationSpreadsheetId,
	})

	if err != nil {
		log.WithError(err).Error("Failed to copy single sheet")
		return DuplicateSpreadsheetResult{}, err
	}

	return DuplicateSpreadsheetResult{
		SpreadsheetId: destinationSpreadsheetId,
	}, nil
}
