package entity

import "time"

type SUserImage struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	UserID    string    `gorm:"type:varchar(36);column:user_id;not null;default:''"`
	TeacherID string    `gorm:"type:varchar(36);column:teacher_id;not null;default:''"`
	StudentID string    `gorm:"type:varchar(36);column:student_id;not null;default:''"`
	ImageID   uint64    `gorm:"column:image_id;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
