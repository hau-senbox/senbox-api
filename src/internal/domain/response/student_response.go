package response

type StudentResponseBase struct {
	StudentID     string              `json:"student_id"`
	StudentName   string              `json:"student_name"`
	Avatar        string              `json:"avatar,omitempty"`
	AvatarURL     string              `json:"avatar_url,omitempty"`
	QrFormProfile string              `json:"qr_form,omitempty"`
	Menus         []ComponentResponse `json:"components,omitempty"`
}
