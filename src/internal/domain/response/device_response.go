package response

import "sen-global-api/internal/domain/value"

type DeviceListSettingResponse struct {
	MacAddress   string `json:"macAddress"`
	DeviceName   string `json:"deviceName"`
	LocationName string `json:"locationName"`
	Note         string `json:"note"`
	DateInstall  string `json:"dateInstall"`
}
type DeviceDetailsAdminResponse struct {
	MacAddress      string `json:"macAddress"`
	DeviceName      string `json:"deviceName"`
	DateInstall     string `json:"dateInstall"`
	SheetId         string `json:"sheetId"`
	SheetName       string `json:"sheetName"`
	SheetLocationId string `json:"sheetLocationId"`
	Location        string `json:"location"`
	LocationId      int64  `json:"locationId"`
	SendEmailTo     string `json:"sendEmailTo"`
	Note            string `json:"note"`
}

type InitDeviceResponseData struct {
	MacAddress string `json:"macAddress"`
}

type DeviceResponse struct {
	Data InitDeviceResponseData `json:"data"`
}

type SendEmailResponseData struct {
	Message bool `json:"message"`
}

type SendEmailResponse struct {
	Data SendEmailResponseData `json:"data"`
}

type AuthorizedDeviceResponseData struct {
	AccessToken  string `json:"accessToken" binding:"required"`
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type AuthorizedDeviceResponse struct {
	Data AuthorizedDeviceResponseData `json:"data"`
}

type DeviceAttributeResponse struct {
	Name string `json:"name"`
}

type DeviceListResponseButton struct {
	ButtonType  string `json:"button_type" binding:"required" enums:"list,scan"`
	ButtonValue string `json:"button_value" binding:"required"`
}

type DeviceListResponseData struct {
	DeviceUUID            string                 `json:"device_uuid"`
	DeviceName            string                 `json:"device_name"`
	Attribute1            string                 `json:"attribute_1"`
	Attribute2            string                 `json:"attribute_2"`
	Attribute3            string                 `json:"attribute_3"`
	InputMode             string                 `json:"input_mode"`
	Status                string                 `json:"status"`
	ProfilePicture        string                 `json:"profile_picture"`
	SpreadsheetUrl        string                 `json:"spreadsheet_url"`
	Message               string                 `json:"message"`
	ButtonUrl             string                 `json:"button_url"`
	ScreenButtonType      value.ScreenButtonType `json:"screen_button_type"`
	ScreenButtonValue     string                 `json:"screen_button_value"`
	AppVersion            string                 `json:"app_version"`
	Note                  string                 `json:"note"`
	UpdatedAt             string                 `json:"updated_at"`
	TeacherSpreadsheetUrl string                 `json:"teacher_spreadsheet_url"`
}

type DeviceListResponse struct {
	Data   []DeviceListResponseData `json:"data"`
	Paging Pagination               `json:"pagination"`
}

type UpdateDeviceResponse struct {
	Data DeviceListResponseData `json:"data"`
}
