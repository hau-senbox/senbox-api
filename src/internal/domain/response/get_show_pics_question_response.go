package response

type GetShowPicsQuestionResponseData struct {
	PhotoURL string `json:"photo_url" binding:"required"`
}

type GetShowPicsQuestionResponse struct {
	Data GetShowPicsQuestionResponseData `json:"data" binding:"required"`
}
