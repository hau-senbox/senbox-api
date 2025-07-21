package response

type GetStudentMenuResponse struct {
	StudentID   string                     `json:"student_id"`
	StudentName string                     `json:"student_name"`
	Components  []ComponentStudentResponse `json:"components"`
}

type ComponentStudentResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Order int    `json:"order"`
	Ishow bool   `json:"is_show"`
}
