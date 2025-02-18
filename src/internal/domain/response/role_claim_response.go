package response

type RoleClaimResponse struct {
	ID         int64  `json:"id"`
	ClaimName  string `json:"claim_name"`
	ClaimValue string `json:"claim_value"`
}

type RoleClaimListResponseData struct {
	ID        int64  `json:"id"`
	ClaimName string `json:"claim_name"`
}

type RoleClaimListResponse struct {
	Data []RoleClaimListResponseData `json:"data"`
}
