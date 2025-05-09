package response

import "time"

type UserListResponseData struct {
	UserID   string `json:"user_id"`
	Fullname string `json:"fullname"`
}

type CreateUserResponseData struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserResponse struct {
	Data CreateUserResponseData `json:"data"`
}
