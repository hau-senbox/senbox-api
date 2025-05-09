package request

type GetDeviceComponentValuesByOrganizationRequest struct {
	ID uint `json:"id" binding:"required"`
}

type GetDeviceComponentValuesByDeviceRequest struct {
	ID uint `json:"id" binding:"required"`
}
