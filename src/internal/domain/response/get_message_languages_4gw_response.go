package response

type GetMessageLanguages4GWResponse struct {
	LangID   uint              `json:"lang_id"`
	Contents map[string]string `json:"contents"`
}
