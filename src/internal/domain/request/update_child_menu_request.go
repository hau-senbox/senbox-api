package request

type UpdateChildMenuRequest struct {
	ChildID     string `json:"child_id" binding:"required"`
	ComponentID string `json:"component_id" binding:"required"`
	IsShow      *bool  `json:"is_show" binding:"required"`
}
