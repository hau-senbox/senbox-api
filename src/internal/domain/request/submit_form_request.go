package request

import "time"

type SubmitFormRequest struct {
	QRCode          string    `json:"qr_code" binding:"required"`
	Answers         []Answer  `json:"answers" binding:"required"`
	OpenedAt        time.Time `json:"opened_at" binding:"required"`
	UserID          string
	ChildID         *string `json:"child_id"`
	StudentCustomID *string `json:"student_custom_id"`
	UserCustomID    *string `json:"user_custom_id"`
}
