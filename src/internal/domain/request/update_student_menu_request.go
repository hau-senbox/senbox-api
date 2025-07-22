package request

type UpdateStudentMenuRequest struct {
	StudentID   string `json:"student_id" binding:"required"`
	ComponentID string `json:"component_id" binding:"required"`
	IsShow      *bool  `json:"is_show" binding:"required"`
}
