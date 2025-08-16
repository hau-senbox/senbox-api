package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SChild struct {
	ID           uuid.UUID `json:"id" gorm:"column:id;type:char(36);primaryKey"`
	ChildName    string    `json:"child_name" gorm:"column:child_name;type:varchar(255);not null"`
	Age          int       `json:"age" gorm:"column:age;not null"`
	ParentID     uuid.UUID `json:"parent_id" gorm:"column:parent_id;type:char(36);not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP;not null"`
	CreatedIndex int       `json:"created_index" gorm:"column:created_index;not null;default:0"`
}

// BeforeCreate hook -> auto generate CreatedIndex = count + 1 (toàn bảng)
func (c *SChild) BeforeCreate(tx *gorm.DB) (err error) {
	var count int64
	if err = tx.Model(&SChild{}).Count(&count).Error; err != nil {
		return err
	}
	c.CreatedIndex = int(count) + 1
	return
}
