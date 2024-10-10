package usecase

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strconv"
)

type UpdateDeviceInfoUseCase struct {
	DeviceRepository  *repository.DeviceRepository
	SettingRepository *repository.SettingRepository
	Reader            *sheet.Reader
	Writer            *sheet.Writer
}

func (c *UpdateDeviceInfoUseCase) Execute(device entity.SDevice, version *string, userInfo3 *string) error {
	setting, err := c.SettingRepository.GetSyncDevicesSettings()
	if err != nil || setting == nil {
		return err
	}

	type ImportSetting struct {
		SpreadSheetUrl string `json:"spreadsheet_url"`
		AutoImport     bool   `json:"auto"`
		Interval       uint8  `json:"interval"`
	}
	var importSetting ImportSetting
	err = json.Unmarshal([]byte(setting.Settings), &importSetting)
	if err != nil {
		log.Error("failed to unmarshal sync devices settings")
		return err
	}

	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(importSetting.SpreadSheetUrl)

	if len(match) < 2 {
		log.Error("failed to parse spreadsheet id from sync devices sheet")
		return err
	}

	spreadsheetId := match[1]

	values, err := c.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     `Devices!K11:AB`,
	})

	if err != nil {
		log.Error("failed to read devices sheet")
		return err
	}

	rowNo, err := c.findFirstRow(device.DeviceId, values, 11)

	if err != nil || rowNo == 0 {
		log.Error("failed to find device id in devices sheet")
		return err
	}

	deviceData := make([][]interface{}, 0)
	deviceData = append(deviceData, []interface{}{version})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{userInfo3})
	_, err = c.Writer.UpdateRange(sheet.WriteRangeParams{
		Range:     "Devices!M" + strconv.Itoa(rowNo),
		Rows:      deviceData,
		Dimension: "COLUMNS",
	}, spreadsheetId)

	if err != nil {
		log.Error("failed to update device info")
		return err
	}

	_ = c.DeviceRepository.UpdateDeviceInfo(device, version, userInfo3)

	//TODO: Write App sheet
	err = c.updateAppSheet(device, version, userInfo3)

	return err
}

func (receiver *UpdateDeviceInfoUseCase) findFirstRow(id string, values [][]interface{}, startRow int) (int, error) {
	rowNo := 0
	for rowindex, row := range values {
		if len(row) > 1 {
			if row[1].(string) == id {
				return rowindex + startRow, nil
			}
		}
	}
	return rowNo, errors.New("Cannot determine row number for device id: " + id)
}

func (c *UpdateDeviceInfoUseCase) updateAppSheet(device entity.SDevice, version *string, info3 *string) error {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(device.ScreenButtonValue)

	if len(match) < 2 {
		log.Error("failed to get spreadsheet id to log accounts")
		return errors.New("failed to get spreadsheet id to log accounts")
	}

	accountSpreadsheetId := match[1]

	//Init Account Sheet
	infoRows := make([][]interface{}, 0)
	infoRows = append(infoRows, []interface{}{info3})
	infoRows = append(infoRows, []interface{}{version})
	accountSheetParams := sheet.WriteRangeParams{
		Range:     "Account!M14",
		Dimension: "ROWS",
		Rows:      infoRows,
	}
	monitor.LogGoogleAPIRequestInitDevice()
	writtenAccountRanges, err := c.Writer.UpdateRange(accountSheetParams, accountSpreadsheetId)
	if err != nil {
		log.Error("failed to write to account sheet")
		return err
	}
	log.Debug(writtenAccountRanges)

	return nil
}
