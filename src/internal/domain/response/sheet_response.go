package response

type ListLocationResponse struct {
	LocationId   int64  `json:"locationId"`
	LocationName string `json:"locationName"`
}
type SettingListResponse struct {
	Uid             string `json:"uid"`
	DeviceName      string `json:"deviceName"`
	LocationId      int64  `json:"locationId"`
	SheetId         string `json:"sheetId"`
	SheetName       string `json:"sheetName"`
	SheetLocationId string `json:"sheetLocationId"`
}
