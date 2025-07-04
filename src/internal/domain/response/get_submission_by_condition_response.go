package response

type GetSubmissionByConditionResponse struct {
	Answer string `json:"answer" binding:"required"`
}
