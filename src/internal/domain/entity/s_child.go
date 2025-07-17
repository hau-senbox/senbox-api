package entity

import (
	"time"

	"github.com/google/uuid"
)

type SChild struct {
	ID        uuid.UUID `json:"id" gorm:"column:id;type:char(36);primaryKey"`
	ChildName string    `json:"child_name" gorm:"column:student_name;type:varchar(255);not null"`
	Age       int       `json:"age" gorm:"column:age;not null"`
	ParentID  uuid.UUID `json:"parent_id" gorm:"column:parent_id;type:char(36);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP;not null"`
}
