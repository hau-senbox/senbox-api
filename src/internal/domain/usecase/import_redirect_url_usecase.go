package usecase

import (
	"fmt"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/pkg/job"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type ImportRedirectUrlsUseCase struct {
	RedirectUrlRepository *repository.RedirectUrlRepository
	SpreadsheetReader     *sheet.Reader
	SpreadsheetWriter     *sheet.Writer
	SettingRepository     *repository.SettingRepository
	TimeMachine           *job.TimeMachine
}

func (receiver *ImportRedirectUrlsUseCase) SyncUrls(req request.ImportRedirectUrlsRequest) error {
	monitor.SendMessageViaTelegram(fmt.Sprintf("[INFO][SYNC] Start sync URLS with interval %d", req.Interval))
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}

	if req.Interval == 0 {
		return nil
	}

	spreadsheetId := match[1]
	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "URL_FORWARD!K12:O",
	})
	if err != nil {
		log.Error(err)
		return err
	}
	for rowNo, row := range values {
		if len(row) >= 5 && cap(row) >= 5 {
			importErr := receiver.RedirectUrlRepository.SaveRedirectUrl(row[1].(string), row[2].(string), row[3].(string), row[4].(string), "", nil)
			if importErr != nil {
				log.Error(importErr)
			} else {
				_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
					Range:     "URL_FORWARD!O" + strconv.Itoa(rowNo+12) + "P",
					Dimension: "ROWS",
					Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05")}},
				}, spreadsheetId)
				if err != nil {
					log.Debug("Row No: ", rowNo)
					log.Error(err)
				}
			}
		}
	}

	return nil
}

func (receiver *ImportRedirectUrlsUseCase) Import(req request.ImportRedirectUrlsRequest) error {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}

	if req.Interval == 0 {
		req.AutoImport = false
	}

	err := receiver.SettingRepository.UpdateUrlSetting(req)
	if err != nil {
		return err
	}

	spreadsheetId := match[1]
	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "URL_FORWARD!K12:T",
	})
	if err != nil {
		log.Error(err)
		return err
	}
	for rowNo, row := range values {
		if len(row) >= 5 && cap(row) >= 5 && row[2].(string) != "" && strings.ToLower(row[4].(string)) == "upload" {
			hint := ""
			if len(row) > 8 {
				hint = row[8].(string)
			}
			var hashPwd *string
			if len(row) > 9 {
				hash := row[9].(string)
				hashPwd = &hash
			}
			importErr := receiver.RedirectUrlRepository.SaveRedirectUrl(row[1].(string), row[2].(string), row[3].(string), row[4].(string), hint, hashPwd)
			if importErr != nil {
				log.Error(importErr)
			} else {
				_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
					Range:     "URL_FORWARD!P" + strconv.Itoa(rowNo+12) + ":Q",
					Dimension: "ROWS",
					Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05")}},
				}, spreadsheetId)
				if err != nil {
					log.Debug("Row No: ", rowNo)
					log.Error(err)
				}
			}
		}
	}

	// if !req.AutoImport {
	// 	receiver.TimeMachine.ScheduleSyncUrls(0)
	// } else {
	// 	receiver.TimeMachine.ScheduleSyncUrls(req.Interval)
	// }

	return nil
}

func (receiver *ImportRedirectUrlsUseCase) ImportPartially(url string, sheetName string) error {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(url)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}
	spreadsheetId := match[1]
	monitor.LogGoogleAPIRequestImportForm()

	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     sheetName + `!K12:T`,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	for rowNo, row := range values {
		if len(row) >= 5 && cap(row) >= 5 {
			hint := ""
			if len(row) > 8 {
				hint = row[8].(string)
			}
			var hashPwd *string
			if len(row) > 9 {
				hash := row[9].(string)
				hashPwd = &hash
			}
			importErr := receiver.RedirectUrlRepository.SaveRedirectUrl(row[1].(string), row[2].(string), row[3].(string), row[4].(string), hint, hashPwd)
			if importErr != nil {
				log.Error(importErr)
			} else {
				_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
					Range:     sheetName + "!P" + strconv.Itoa(rowNo+12) + ":Q",
					Dimension: "ROWS",
					Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05")}},
				}, spreadsheetId)
				if err != nil {
					log.Debug("Row No: ", rowNo)
					log.Error(err)
				}
			}
		}
	}

	return nil
}
