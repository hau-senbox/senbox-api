package request

type GenerateOwnerCodeRequest struct {
	OwnerID      string `json:"owner_id"`
	CreatedIndex int    `json:"created_index"`
}
