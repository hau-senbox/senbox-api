package response

type PdfResponse struct {
	PdfName        string `json:"pdf_name"`
	Key            string `json:"key"`
	OrganizationID uint64 `json:"organization_id"`
	Url            string `json:"url"`
	Extension      string `json:"extension"`
}
