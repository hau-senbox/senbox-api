package response

type UserEntityResponse struct {
	ID           string   `json:"id"`
	Username     string   `json:"username"`
	Nickname     string   `json:"nickname"`
	Fullname     string   `json:"fullname"`
	Phone        string   `json:"phone"`
	Email        string   `json:"email"`
	Dob          string   `json:"dob"`
	QRLogin      string   `json:"qr_login"`
	Avatar       string   `json:"avatar"`
	AvatarURL    string   `json:"avatar_url"`
	IsBlocked    bool     `json:"is_blocked"`
	BlockedAt    string   `json:"blocked_at"`
	Organization []string `json:"organizations"`
	CreatedAt    string   `json:"created_at"`

	Roles   *[]RoleListResponseData `json:"roles"`
	Devices *[]string               `json:"devices"`
}

type UserEntityResponseV2 struct {
	ID           string   `json:"id"`
	Username     string   `json:"username"`
	Nickname     string   `json:"nickname"`
	Fullname     string   `json:"fullname"`
	Phone        string   `json:"phone"`
	Email        string   `json:"email"`
	Dob          string   `json:"dob"`
	QRLogin      string   `json:"qr_login"`
	Avatar       string   `json:"avatar"`
	AvatarURL    string   `json:"avatar_url"`
	IsBlocked    bool     `json:"is_blocked"`
	BlockedAt    string   `json:"blocked_at"`
	Organization []string `json:"organizations"`
	CreatedAt    string   `json:"created_at"`

	Roles   *[]RoleListResponseData `json:"roles"`
	Devices *[]string               `json:"devices"`
}

type UserEntityResponseData struct {
	ID        string   `json:"id"`
	Username  string   `json:"username"`
	Nickname  string   `json:"nickname"`
	Avatar    string   `json:"avatar"`
	AvatarURL string   `json:"avatar_url"`
	Roles     []string `json:"roles"`
}
