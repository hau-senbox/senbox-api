package request

type UploadSectionUserMenuRequest UserSectionMenuItem

type UserSectionMenuItem struct {
	LanguageID         uint                         `json:"language_id" binding:"required"`
	UserID             string                       `json:"user_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
