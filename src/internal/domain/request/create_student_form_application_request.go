package request

type CreateStudentFormApplicationRequest struct {
	StudentName    string `json:"student_name" binding:"required"`
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
	ChildID        string `json:"child_id" binding:"required"`
}
