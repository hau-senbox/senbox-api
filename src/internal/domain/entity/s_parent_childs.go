package entity

type SParentChilds struct {
	ParentID string `json:"parent_id" gorm:"type:char(36);not null"`
	ChildID  string `json:"child_id" gorm:"type:char(36);not null"`
}
