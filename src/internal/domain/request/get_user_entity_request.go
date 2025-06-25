package request

type GetUserEntityByIDRequest struct {
	ID string `json:"id" binding:"required"`
}

type GetUserEntityByUsernameRequest struct {
	Username string `json:"username" binding:"required"`
}
