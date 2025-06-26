package entity

import "github.com/google/uuid"

type SUserParentChild struct {
	ParentID uuid.UUID   `gorm:"column:parent_id;primary_key"`
	Parent   SUserEntity `gorm:"foreignKey:ParentID;references:id;constraint:OnDelete:CASCADE;"`
	ChildID  uuid.UUID   `gorm:"column:child_id;primary_key"`
	Child    SUserEntity `gorm:"foreignKey:ChildID;references:id;constraint:OnDelete:CASCADE;"`
}
