package response

import "time"

type StudentBlockSettingResponse struct {
	ID              int       `json:"id"`
	StudentID       string    `json:"student_id"`
	IsDeactive      bool      `json:"is_deactive"`
	IsViewMessage   bool      `json:"is_view_message"`
	MessageBox      string    `json:"message_box"`
	MessageDeactive string    `json:"message_deactive"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
