package request

type AddCustomID2UserRequest struct {
	UserID   string `json:"user_id"`
	CustomID string `json:"custom_id"`
}
