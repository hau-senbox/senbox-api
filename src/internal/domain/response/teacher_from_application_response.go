package response

type TeacherFormApplicationResponse struct {
	ID         int64  `json:"id"`
	Status     string `json:"status"`
	ApprovedAt string `json:"approved_at"`
	CreatedAt  string `json:"created_at"`
	UserID     string `json:"user_id"`
}
