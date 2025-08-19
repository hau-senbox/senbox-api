package entity

import (
	"sen-global-api/internal/domain/value"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserImages struct {
	ID        uuid.UUID       `json:"id" gorm:"type:char(36);primaryKey"`
	OwnerID   uuid.UUID       `json:"owner_id" gorm:"type:char(36);not null"`
	OwnerRole value.OwnerRole `json:"owner_role" gorm:"type:varchar(50);not null"`
	ImageID   uint64          `json:"image_id" gorm:"not null"`
	Index     int             `json:"index" gorm:"not null;default:0"`
	IsMain    bool            `json:"is_main" gorm:"not null;default:false"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

// Generate UUID before insert
func (u *UserImages) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
