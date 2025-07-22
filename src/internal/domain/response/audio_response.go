package response

type AudioResponse struct {
	AudioName string `json:"video_name"`
	Key       string `json:"key"`
	Extension string `json:"extension"`
	Url       string `json:"url"`
}
