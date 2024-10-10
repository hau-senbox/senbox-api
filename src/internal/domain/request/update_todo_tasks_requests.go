package request

type UpdateTask struct {
	Name      string `json:"name" binding:"required"`
	DueDate   string `json:"due_date" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Selection string `json:"selection" binding:"required"`
	Selected  string `json:"selected" binding:"required"`
}

type UpdateToDoTasksRequest struct {
	QRCode string       `json:"qr_code" binding:"required"`
	Name   string       `json:"name"`
	Tasks  []UpdateTask `json:"tasks" binding:"required"`
}
