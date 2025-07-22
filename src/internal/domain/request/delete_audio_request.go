package request

type DeleteAudioRequest struct {
	Key string `json:"key" binding:"required"`
}
