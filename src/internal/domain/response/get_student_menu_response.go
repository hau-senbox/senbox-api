package response

type GetStudentMenuResponse struct {
	StudentID   string              `json:"student_id"`
	StudentName string              `json:"student_name"`
	MenuIconKey string              `json:"menu_icon_key"`
	Components  []ComponentResponse `json:"components"`
}
