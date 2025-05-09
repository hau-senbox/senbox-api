package response

type FunctionClaimResponse struct {
	ID           int64                                     `json:"id"`
	FunctionName string                                    `json:"function_name"`
	Permissions  []FunctionClaimPermissionListResponseData `json:"permissions"`
}

type FunctionClaimListResponseData struct {
	ID           int64                                     `json:"id"`
	FunctionName string                                    `json:"function_name"`
	Permissions  []FunctionClaimPermissionListResponseData `json:"permissions"`
}
