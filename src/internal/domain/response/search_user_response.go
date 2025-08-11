package response

type SearchUserResponse struct {
	Users    []UserResponse     `json:"users"`
	Children []ChildrenResponse `json:"children"`
	Students []StudentResponse  `json:"students"`
	Teachers []TeacherResponse  `json:"teachers"`
	Staffs   []StaffResponse    `json:"staffs"`
	Parents  []ParentResponse   `json:"parents"`
}

type UserResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	AvatarURL  string `json:"avatar_url"`
	IsDeactive bool   `json:"is_deactive"`
}

type ChildrenResponse struct {
	ChildID   string `json:"id"`
	ChildName string `json:"nickname"`
}

type StudentResponse struct {
	StudentID   string `json:"id"`
	StudentName string `json:"nickname"`
}

type TeacherResponse struct {
	TeacherID   string `json:"id"`
	TeacherName string `json:"nickname"`
	IsDeactive  bool   `json:"is_deactive"`
}

type StaffResponse struct {
	StaffID    string `json:"id"`
	StaffName  string `json:"nickname"`
	IsDeactive bool   `json:"is_deactive"`
}

type ParentResponse struct {
	ParentID   string `json:"id"`
	ParentName string `json:"nickname"`
	IsDeactive bool   `json:"is_deactive"`
}
