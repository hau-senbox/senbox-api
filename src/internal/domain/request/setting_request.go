package request

type SettingEntryRequest struct {
	Uid             string `json:"uid"`
	LocationID      int64  `json:"locationID"`
	SheetID         string `json:"sheetID"`
	SheetName       string `json:"sheetName"`
	SheetLocationID string `json:"sheetLocationID"`
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
	SheetID         string `json:"sheetID"`
	SheetName       string `json:"sheetName"`
	SheetLocationID string `json:"sheetLocationID"`
	LocationID      int64  `json:"locationID"`
	SendEmailTo     string `json:"sendEmailTo"`
	Note            string `json:"note"`
}
