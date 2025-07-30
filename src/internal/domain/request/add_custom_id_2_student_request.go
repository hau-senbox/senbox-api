package request

type AddCustomId2StudentRequest struct {
	StudentID string `json:"student_id"`
	CustomID  string `json:"custom_id"`
}
