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
	ID           string `json:"id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	IsDeactive   bool   `json:"is_deactive"`
	Avatar       Avatar `json:"avatar"`
	CreatedIndex int    `json:"created_index"`
}

type ChildrenResponse struct {
	ChildID      string `json:"id"`
	ChildName    string `json:"nickname"`
	CreatedIndex int    `json:"created_index"`
	Avatar       Avatar `json:"avatar"`
}

type StudentResponse struct {
	StudentID    string `json:"id"`
	StudentName  string `json:"nickname"`
	IsDeactive   bool   `json:"is_deactive"`
	CreatedIndex int    `json:"created_index"`
	Avatar       Avatar `json:"avatar"`
}

type TeacherResponse struct {
	TeacherID    string `json:"id"`
	TeacherName  string `json:"nickname"`
	IsDeactive   bool   `json:"is_deactive"`
	CreatedIndex int    `json:"created_index"`
	Avatar       Avatar `json:"avatar"`
}

type StaffResponse struct {
	StaffID      string `json:"id"`
	StaffName    string `json:"nickname"`
	IsDeactive   bool   `json:"is_deactive"`
	CreatedIndex int    `json:"created_index"`
	Avatar       Avatar `json:"avatar"`
}

type ParentResponse struct {
	ParentID     string `json:"id"`
	ParentName   string `json:"nickname"`
	IsDeactive   bool   `json:"is_deactive"`
	CreatedIndex int    `json:"created_index"`
	Avatar       Avatar `json:"avatar"`
}
