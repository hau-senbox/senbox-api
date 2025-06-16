package request

type GetDeviceComponentValuesByOrganizationRequest struct {
	ID string `json:"id" binding:"required"`
}

type GetDeviceComponentValuesByDeviceRequest struct {
	ID string `json:"id" binding:"required"`
}
