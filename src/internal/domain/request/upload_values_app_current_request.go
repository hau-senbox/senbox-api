package request

type UploadValuesAppCurrentRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	Value1   string `json:"value1" binding:"required"`
	Value2   string `json:"value2" binding:"required"`
	Value3   string `json:"value3" binding:"required"`
	ImageKey string `json:"image_key"`
}
