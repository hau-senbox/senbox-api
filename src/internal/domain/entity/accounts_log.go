package entity

import (
	"sen-global-api/internal/domain/value"
	"time"
)

type AccountsLog struct {
	ID             int64                 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Type           value.AccountsLogType `gorm:"column:type;not null"`
	UserID         string                `gorm:"column:user_id;type:varchar(255);not null default ''"`
	OrganizationID string                `gorm:"column:organization_id;type:varchar(255);not null default ''"`
	DeviceID       string                `gorm:"column:device_id;type:varchar(255);not null default ''"`
	Created        time.Time             `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
