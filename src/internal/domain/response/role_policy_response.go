package response

type RolePolicyResponse struct {
	ID          int64  `json:"id"`
	PolicyName  string `json:"policy_name"`
	Description string `json:"description"`
}

type RolePolicyListResponseData struct {
	ID         int64  `json:"id"`
	PolicyName string `json:"policy_name"`
}

type RolePolicyListResponse struct {
	Data []RolePolicyListResponseData `json:"data"`
}
