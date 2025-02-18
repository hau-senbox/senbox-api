package request

type UpdateUserGuardianRequest struct {
	UserId    string   `json:"user_id" binding:"required"`
	Guardians []string `json:"guardians" binding:"required"`
}
