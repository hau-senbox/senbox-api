package request

type UploadDeviceInfoRequest struct {
	DeviceName string `json:"device_name" binding:"required"`
}
