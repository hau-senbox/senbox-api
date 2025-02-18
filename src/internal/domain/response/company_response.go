package response

type CompanyResponse struct {
	ID          int64  `json:"id"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Description string `json:"description"`
}
