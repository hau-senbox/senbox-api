package response

type GetMessageLanguages4GWResponse struct {
	LangID   uint              `json:"language_id"`
	Contents map[string]string `json:"contents"`
}
