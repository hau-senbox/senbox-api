package request

type SettingEntryRequest struct {
	Uid             string `json:"uid"`
	LocationId      int64  `json:"locationId"`
	SheetId         string `json:"sheetId"`
	SheetName       string `json:"sheetName"`
	SheetLocationId string `json:"sheetLocationId"`
	DeviceName      string `json:"deviceName"`
	Mac             string `json:"mac"`
}
type InitDeviceNameRequest struct {
	DeviceName   string `json:"deviceName"`
	MacAddress   string `json:"macAddress"`
	Lat          string `json:"lat"`
	Lot          string `json:"lot"`
	LocationName string `json:"locationName"`
}
type GetSettingDetailsRequest struct {
	MacAddress string `json:"macAddress"`
}
type DeleteDeviceRequest struct {
	MacAddress string `json:"macAddress"`
}
type AdminDeviceUpdateRequest struct {
	MacAddress      string `json:"macAddress"`
	DeviceName      string `json:"deviceName"`
	SheetId         string `json:"sheetId"`
	SheetName       string `json:"sheetName"`
	SheetLocationId string `json:"sheetLocationId"`
	LocationId      int64  `json:"locationId"`
	SendEmailTo     string `json:"sendEmailTo"`
	Note            string `json:"note"`
}
