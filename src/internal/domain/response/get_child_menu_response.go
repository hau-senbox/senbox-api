package response

type GetChildMenuResponse struct {
	ChildID    string              `json:"child_id"`
	ChildName  string              `json:"child_name"`
	Components []ComponentResponse `json:"components"`
}
