package request

type GetAnswerByKeyAndDB struct {
	Key string `json:"key" binding:"required"`
	DB  string `json:"db" binding:"required"`
}
