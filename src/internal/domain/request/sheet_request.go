package request

type SheetEntryRequest struct {
	SheetId        string `json:"sheetId" binding:"required"`
	SheetName      string `json:"sheetName" binding:"required"`
	DeviceName     string `json:"deviceName" binding:"required"`
	Count          int64  `json:"count" binding:"required"`
	Lat            string `json:"lat" binding:"required"`
	Lot            string `json:"lot" binding:"required"`
	TimeUp         string `json:"timeUp" binding:"required"`
	ValueOfBarcode string `json:"valueOfBarcode" binding:"required"`
	Note           string `json:"note" binding:"required"`
	Location       string `json:"location" binding:"required"`
	LocationName   string `json:"locationName" binding:"required"`
}

type LocationSheetRequest struct {
	SheetId string `json:"sheetId" binding:"required"`
}
