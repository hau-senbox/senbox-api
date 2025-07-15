package response

import "encoding/json"

type GetCommonMenuResponse struct {
	Component []ComponentResponse `json:"components"`
}

type ComponentResponse struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Type  string          `json:"type"`
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
	Order int             `json:"order"`
}
