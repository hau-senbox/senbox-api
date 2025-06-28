package request

type CreateUserMenuRequest struct {
	UserID     string                       `json:"user_id" binding:"required"`
	Components []CreateMenuComponentRequest `json:"components" binding:"required"`
}
