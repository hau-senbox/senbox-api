package response

import "time"

type UserBlockSettingResponse struct {
	ID              int       `json:"id"`
	UserID          string    `json:"user_id"`
	IsDeactive      bool      `json:"is_deactive"`
	IsViewMessage   bool      `json:"is_view_message"`
	MessageBox      string    `json:"message_box"`
	MessageDeactive string    `json:"message_deactive"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
