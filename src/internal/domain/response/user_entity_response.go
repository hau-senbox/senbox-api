package response

type UserEntityResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Dob       string `json:"dob"`
	Company   string `json:"company_name"`
	CreatedAt string `json:"created_at"`

	UserConfig   *UserConfigResponse           `json:"user_config"`
	Roles        *[]RoleListResponseData       `json:"roles"`
	RolePolicies *[]RolePolicyListResponseData `json:"role_policies"`
	Guardians    *[]UserEntityResponseData     `json:"guardians"`
	Devices      *[]string                     `json:"devices"`
}

type UserEntityResponseData struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type UserEntityDataResponse struct {
	Data []UserEntityResponseData `json:"data"`
}
