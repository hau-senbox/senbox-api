package request

type UpdateCompanyRequest struct {
	ID          int64  `json:"id" binding:"required"`
	CompanyName string `json:"company_name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Description string `json:"description" binding:"required"`
}
