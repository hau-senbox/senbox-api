package request

type UpdateTeacherMenuRequest struct {
	TeacherID   string `json:"teacher_id" binding:"required"`
	ComponentID string `json:"component_id" binding:"required"`
	IsShow      *bool  `json:"is_show" binding:"required"`
}
