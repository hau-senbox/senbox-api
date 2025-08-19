package response

type GetTeacherMenuResponse struct {
	TeacherID   string              `json:"teacher_id"`
	TeacherName string              `json:"teacher_name"`
	MenuIconKey string              `json:"menu_icon_key"`
	Components  []ComponentResponse `json:"components"`
}
