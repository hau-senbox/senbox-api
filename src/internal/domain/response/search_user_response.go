package response

type SearchUserResponse struct {
	Users    []UserResponse     `json:"users"`
	Children []ChildrenResponse `json:"children"`
	Students []StudentResponse  `json:"students"`
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
	ChildName string `json:"username"`
}

type StudentResponse struct {
	StudentID   string `json:"id"`
	StudentName string `json:"username"`
}
