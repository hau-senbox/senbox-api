package request

type TakeNoteRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	Note     string `json:"note" binding:"required"`
}
