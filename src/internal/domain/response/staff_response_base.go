package response

type StaffResponseBase struct {
	StaffID       string              `json:"id"`
	UserID        string              `json:"user_id"`
	StaffName     string              `json:"name"`
	Avatar        string              `json:"avatar,omitempty"`
	AvatarURL     string              `json:"avatar_url,omitempty"`
	QrFormProfile string              `json:"qr_form,omitempty"`
	Menus         []ComponentResponse `json:"components"`
	IsUserBlock   bool                `json:"is_block"`
}
