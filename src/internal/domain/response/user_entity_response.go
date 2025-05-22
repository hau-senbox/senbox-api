package response

type UserEntityResponse struct {
	ID           string   `json:"id"`
	Username     string   `json:"username"`
	Nickname     string   `json:"nickname"`
	Fullname     string   `json:"fullname"`
	Phone        string   `json:"phone"`
	Email        string   `json:"email"`
	Dob          string   `json:"dob"`
	IsBlocked    bool     `json:"is_blocked"`
	BlockedAt    string   `json:"blocked_at"`
	Organization []string `json:"organizations"`
	CreatedAt    string   `json:"created_at"`

	Roles     *[]RoleListResponseData   `json:"roles"`
	Guardians *[]UserEntityResponseData `json:"guardians"`
	Devices   *[]string                 `json:"devices"`
}

type UserEntityResponseData struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Roles    []string `json:"roles"`
}
