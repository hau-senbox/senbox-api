package request

type GetDeviceComponentValuesByCompanyRequest struct {
	ID uint `json:"id" binding:"required"`
}

type GetDeviceComponentValuesByDeviceRequest struct {
	ID uint `json:"id" binding:"required"`
}
