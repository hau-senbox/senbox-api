package response

type Avatar struct {
	ImageID  uint64 `json:"image_id"`
	ImageKey string `json:"image_key"`
	Index    int    `json:"index"`
	IsMain   bool   `json:"is_main"`
}
