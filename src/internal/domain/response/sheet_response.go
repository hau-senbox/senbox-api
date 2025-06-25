package response

type ListLocationResponse struct {
	LocationID   int64  `json:"locationID"`
	LocationName string `json:"locationName"`
}
type SettingListResponse struct {
	Uid             string `json:"uid"`
	DeviceName      string `json:"deviceName"`
	LocationID      int64  `json:"locationID"`
	SheetID         string `json:"sheetID"`
	SheetName       string `json:"sheetName"`
	SheetLocationID string `json:"sheetLocationID"`
}
