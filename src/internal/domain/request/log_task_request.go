package request

type LogTaskType string

const (
	LogTaskType_Create  = "created"
	LogTaskType_Update  = "updated"
	LogTaskType_Deleted = "deleted"
)

type LogTaskRequest struct {
	ToDoID  string      `json:"qr_code" binding:"required"`
	Name    string      `json:"name" binding:"required"`
	DueDate string      `json:"due_date" binding:"required"`
	Value   string      `json:"value"`
	LogType LogTaskType `json:"type" binding:"required"`
}
