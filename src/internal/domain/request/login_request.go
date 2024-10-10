package request

type LoginInputReq struct {
	LoginId  string `json:"loginId" binding:"required"`
	Password string `json:"password" binding:"required"`
}
