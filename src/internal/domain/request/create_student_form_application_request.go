package request

type CreateStudentFormApplicationRequest struct {
	StudentName    string `json:"student_name" binding:"required"`
	UserID         string `json:"user_id"`
	OrganizationID int64  `json:"organization_id"`
}
