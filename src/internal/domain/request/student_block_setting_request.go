package request

type StudentBlockSettingRequest struct {
	StudentID       string `json:"student_id"`
	IsDeactive      *bool  `json:"is_deactive" binding:"required"`
	IsViewMessage   *bool  `json:"is_view_message" binding:"required"`
	MessageBox      string `json:"message_box"`
	MessageDeactive string `json:"message_deactive"`
}
