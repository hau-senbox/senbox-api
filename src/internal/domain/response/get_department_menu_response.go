package response

type GetDepartmentMenuResponse struct {
	DepartmentID   string              `json:"department_id,omitempty"`
	DepartmentName string              `json:"department_name,omitempty"`
	Components     []ComponentResponse `json:"components"`
}
