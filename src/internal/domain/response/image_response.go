package response

type ImageResponse struct {
	ImageName string `json:"image_name"`
	Key       string `json:"key"`
	Extension string `json:"extension"`
	Url       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}
