package response

type GetTimeTableItem struct {
	StartAt      string `json:"start_at" binding:"required"`
	EndAt        string `json:"end_at" binding:"required"`
	Color        string `json:"color" binding:"required"`
	Message      string `json:"message" binding:"required"`
	Notification string `json:"notification" binding:"required"`
	Picture      string `json:"picture" binding:"required"`
	Link         string `json:"link" binding:"required"`
}

type GetTimeTableResponseData struct {
	GeneralMessage       string             `json:"general_message" binding:"required"`
	NumberOfItemsPerTime int                `json:"number_of_items_per_time" binding:"required"`
	Times                []GetTimeTableItem `json:"times" binding:"required"`
}

type GetTimeTableResponse struct {
	Data GetTimeTableResponseData `json:"data" bing:"required"`
}
