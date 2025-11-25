package response

import "time"

type GetValuesAppResponse struct {
	Value1    string    `json:"value1"`
	Value2    string    `json:"value2"`
	Value3    string    `json:"value3"`
	ImageKey  string    `json:"image_key"`
	ImageUrl  *string   `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}
