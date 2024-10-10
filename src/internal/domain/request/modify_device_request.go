package request

type ModifyDeviceRequest struct {
	DeviceName               *string `json:"device_name"`
	Note                     *string `json:"note"`
	Message                  *string `json:"message"`
	Status                   *string `json:"status"`
	AppSettingSpreadsheetUrl *string `json:"app_setting_spreadsheet_url"`
	OutputSpreadsheetUrl     *string `json:"output_spreadsheet_url"`
}
