package request

type GetCodeCountingsRequest struct {
	Keyword string `form:"keyword"`
	PageNo  int    `form:"page" default:"1"`
	PerPage int    `form:"limit" default:"12"`
}
