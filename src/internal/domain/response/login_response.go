package response

import "time"

type LoginResponseData struct {
	UserId   string    `json:"user_id"`
	Username string    `json:"username"`
	Token    string    `json:"token"`
	Expired  time.Time `json:"expired"`
}

type LoginResponse struct {
	Data LoginResponseData `json:"data"`
}
