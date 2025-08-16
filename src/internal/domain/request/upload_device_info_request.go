package request

type UploadDeviceInfoRequest struct {
	DeviceNickName string `json:"device_nick_name" binding:"required"`
}
