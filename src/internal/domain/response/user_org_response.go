package response

type GetUserOrgInfoResponse struct {
	UserNickName string `json:"user_nick_name"`
	IsManager    bool   `json:"is_manager"`
}

type GetOrgManagerInfoResponse struct {
	UserID       string `json:"user_id"`
	UserNickName string `json:"user_nick_name"`
	IsManager    bool   `json:"is_manager"`
}
