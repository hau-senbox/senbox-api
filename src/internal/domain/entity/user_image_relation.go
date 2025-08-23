package entity

import (
	"time"
)

// UserImageRelation: chỉ định một người khác có liên quan đến hình ảnh (ngoài Owner chính).
type UserImageRelation struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserImageID string    `json:"user_image_id" gorm:"type:char(36);not null;index"`
	RelatedID   string    `json:"related_id" gorm:"type:char(36);not null"`
	RelatedRole string    `json:"related_role" gorm:"type:varchar(50);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
