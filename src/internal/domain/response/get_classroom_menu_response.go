package response

type GetClassroomMenuResponse struct {
	ClassroomID   string              `json:"classroom_id,omitempty"`
	ClassroomName string              `json:"classroom_name,omitempty"`
	Components    []ComponentResponse `json:"components"`
}
