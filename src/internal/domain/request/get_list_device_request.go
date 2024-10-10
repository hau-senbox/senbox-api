package request

type GetListDeviceRequest struct {
	Page    int    `form:"page"`
	Keyword string `form:"keyword"`
	Limit   int    `form:"limit"`
}
