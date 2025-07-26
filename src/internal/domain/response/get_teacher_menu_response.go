package response

type GetTeacherMenuResponse struct {
	TeacherID   string              `json:"teacher_id"`
	TeacherName string              `json:"teacher_name"`
	Components  []ComponentResponse `json:"components"`
}
