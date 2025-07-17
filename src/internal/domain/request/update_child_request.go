package request

type UpdateChildRequest struct {
	ID        string `json:"id" binding:"required"`
	ChildName string `json:"child_name" binding:"required"`
	Age       int    `json:"age" binding:"required"`
}
