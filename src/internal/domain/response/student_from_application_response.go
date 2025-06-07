package response

type StudentFormApplicationResponse struct {
	ID          int64  `json:"id"`
	StudentName string `json:"student_name"`
	Status      string `json:"status"`
	ApprovedAt  string `json:"approved_at"`
	CreatedAt   string `json:"created_at"`
	UserID      string `json:"user_id"`
}
