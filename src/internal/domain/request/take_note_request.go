package request

type TakeNoteRequest struct {
	DeviceId string `json:"device_id" binding:"required"`
	Note     string `json:"note" binding:"required"`
}
