package response

type GetStaffMenuResponse struct {
	StaffID    string              `json:"staff_id"`
	StaffName  string              `json:"staff_name"`
	Components []ComponentResponse `json:"components"`
}
