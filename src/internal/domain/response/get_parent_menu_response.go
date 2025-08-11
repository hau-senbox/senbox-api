package response

type GetParentMenuResponse struct {
	ParentID   string              `json:"parent_id"`
	ParentName string              `json:"parent_name"`
	Components []ComponentResponse `json:"components"`
}
