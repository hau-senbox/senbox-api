package response

type StudentResponseBase struct {
	StudentID     string              `json:"student_id"`
	StudentName   string              `json:"student_name"`
	Avatar        string              `json:"avatar"`
	AvatarURL     string              `json:"avatar_url"`
	QrFormProfile string              `json:"qr_form"`
	Menus         []ComponentResponse `json:"components"`
}
