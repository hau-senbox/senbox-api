package response

import "time"

type LoginResponseData struct {
	UserId   string    `json:"userId"`
	UserName string    `json:"userName"`
	Token    string    `json:"token"`
	Expired  time.Time `json:"expired"`
}

type LoginResponse struct {
	Data LoginResponseData `json:"data"`
}
