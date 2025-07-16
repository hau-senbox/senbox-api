package response

type GetCommonMenuResponse struct {
	Component []ComponentResponse `json:"components"`
}

type ComponentResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Order int    `json:"order"`
}
