package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SParent struct {
	ID           uuid.UUID `gorm:"column:id;type:char(36);primary_key"`
	UserID       string    `gorm:"column:user_id;type:char(36);not null"`
	CreatedIndex int       `gorm:"column:created_index;not null;default:0"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime:milli;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime:milli;not null"`
	ParentName   string    `gorm:"-" json:"parent_name"` // not map db
}

func (parent *SParent) BeforeCreate(tx *gorm.DB) (err error) {

	// Tính CreatedIndex = MAX(created_index) + 1 theo UserID
	var maxIndex int
	if err := tx.Model(&SParent{}).
		Select("COALESCE(MAX(created_index), 0)").
		Scan(&maxIndex).Error; err != nil {
		return err
	}

	parent.CreatedIndex = maxIndex + 1

	return nil
}
