package response

import "time"

type UserListResponseData struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type UserListResponse struct {
	Data []UserListResponseData `json:"data"`
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
