package response

type PdfResponse struct {
	PdfName   string `json:"pdf_name"`
	Key       string `json:"key"`
	Url       string `json:"url"`
	Extension string `json:"extension"`
}
