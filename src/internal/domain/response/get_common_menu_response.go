package response

type GetCommonMenuResponse struct {
	ChildID    string                `json:"child_id"`
	Components []ComponentCommonMenu `json:"components"`
}

type ComponentResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Order int    `json:"order"`
}

type ComponentCommonMenu struct {
	ChildID    string              `json:"child_id"`
	Components []ComponentResponse `json:"components"`
}
