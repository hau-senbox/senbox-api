package response

type SearchUserResponse struct {
	Users    []UserResponse     `json:"users"`
	Children []ChildrenResponse `json:"children"`
	Students []StudentResponse  `json:"students"`
	Teachers []TeacherResponse  `json:"teachers"`
	Staffs   []StaffResponse    `json:"staffs"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	AvatarURL string `json:"avatar_url"`
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
}

type StaffResponse struct {
	StaffID   string `json:"id"`
	StaffName string `json:"nickname"`
}
