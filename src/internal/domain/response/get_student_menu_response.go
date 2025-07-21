package response

type GetStudentMenuResponse struct {
	StudentID   string              `json:"student_id"`
	StudentName string              `json:"student_name"`
	Components  []ComponentResponse `json:"components"`
}
