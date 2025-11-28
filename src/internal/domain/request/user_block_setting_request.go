package request

type UserBlockSettingRequest struct {
	UserID          string `json:"user_id"`
	IsDeactive      *bool  `json:"is_deactive" binding:"required"`
	IsViewMessage   *bool  `json:"is_view_message" binding:"required"`
	MessageBox      string `json:"message_box"`
	MessageDeactive string `json:"message_deactive"`
	IsNeedToUpdate  bool   `json:"is_need_to_update"`
}
