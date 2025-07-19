package response

type GetChildMenuResponse struct {
	ChildID    string                   `json:"child_id"`
	ChildName  string                   `json:"child_name"`
	Components []ComponentChildResponse `json:"components"`
}

type ComponentChildResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Order int    `json:"order"`
	Ishow bool   `json:"is_show"`
}
