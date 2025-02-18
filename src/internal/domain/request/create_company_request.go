package request

type CreateCompanyRequest struct {
	CompanyName string `json:"company_name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Description string `json:"description" binding:"required"`
}
