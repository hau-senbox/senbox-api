package response

type StudentResponseBase struct {
	StudentID     string              `json:"id"`
	StudentName   string              `json:"name"`
	Avatar        string              `json:"avatar,omitempty"`
	AvatarURL     string              `json:"avatar_url,omitempty"`
	QrFormProfile string              `json:"qr_form,omitempty"`
	Menus         []ComponentResponse `json:"components"`
	CustomID      string              `json:"custom_id"`
}
