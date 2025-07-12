package response

type ChartDataResponse struct {
	X string `json:"x" bindding:"required"`
	Y string `json:"y" binding:"required"`
}
