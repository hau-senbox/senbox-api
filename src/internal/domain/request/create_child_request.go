package request

type CreateChildRequest struct {
	ChildName string `json:"child_name" binding:"required"`
	Age       int    `json:"age" binding:"required"`
}
