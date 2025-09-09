package request

type UpdateReLoginWebRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	ReLogin *bool  `json:"re_login" binding:"required"`
}
