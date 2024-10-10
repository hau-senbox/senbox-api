package request

type TakeNoteRequest struct {
	Note string `json:"note" binding:"required"`
}
