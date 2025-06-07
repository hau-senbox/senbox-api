package request

type UploadSuperAdminMenuRequest struct {
	Top    []CreateMenuComponentRequest `json:"top" binding:"required"`
	Bottom []CreateMenuComponentRequest `json:"bottom" binding:"required"`
}
