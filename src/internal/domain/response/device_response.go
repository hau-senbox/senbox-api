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

type DeviceAttributeResponse struct {
	Name string `json:"name"`
}

type DeviceListResponseButton struct {
	ButtonType  string `json:"button_type" binding:"required" enums:"list,scan"`
	ButtonValue string `json:"button_value" binding:"required"`
}

type DeviceResponseData struct {
	DeviceUUID        string                 `json:"device_uuid"`
	DeviceName        string                 `json:"device_name"`
	InputMode         string                 `json:"input_mode"`
	Status            string                 `json:"status"`
	DeactivateMessage string                 `json:"deactivate_message"`
	ButtonUrl         string                 `json:"button_url"`
	ScreenButtonType  value.ScreenButtonType `json:"screen_button_type"`
	AppVersion        string                 `json:"app_version"`
	Note              string                 `json:"note"`
	UpdatedAt         string                 `json:"updated_at"`
}

type DeviceResponseDataV2 struct {
	Id                string `json:"id"`
	DeviceName        string `json:"device_name"`
	InputMode         string `json:"input_mode"`
	Status            string `json:"status"`
	DeactivateMessage string `json:"deactivate_message"`
	ButtonUrl         string `json:"button_url"`
	AppVersion        string `json:"app_version"`
	Note              string `json:"note"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type DeviceListResponse struct {
	Data   []DeviceResponseData `json:"data"`
	Paging Pagination           `json:"pagination"`
}

type DeviceListResponseV2 struct {
	Data   []DeviceResponseDataV2 `json:"data"`
	Paging Pagination             `json:"pagination"`
}

type UpdateDeviceResponse struct {
	Data DeviceResponseData `json:"data"`
}

type DeviceResponseV2 struct {
	ID         string `json:"id"`
	DeviceName string `json:"device_name"`
}
