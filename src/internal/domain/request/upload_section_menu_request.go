package request

type UploadSectionMenuRequest struct {
	Components []CreateMenuComponentRequest `json:"components" binding:"required"`
}
