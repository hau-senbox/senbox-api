package entity

import "time"

type ValuesAppHistories struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	DeviceID  string    `gorm:"column:device_id;type:varchar(100);not null;default:''"`
	Value1    string    `gorm:"column:value1;type:varchar(50);not null;default:''"`
	Value2    string    `gorm:"column:value2;type:varchar(50);not null;default:''"`
	Value3    string    `gorm:"column:value3;type:varchar(50);not null;default:''"`
	ImageKey  string    `gorm:"column:image_key;type:varchar(100);not null;default:''"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP;not null"`
}
