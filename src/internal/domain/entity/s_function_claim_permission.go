package entity

import (
	"time"
)

type SFunctionClaimPermission struct {
	ID              int64          `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	PermissionName  string         `gorm:"type:varchar(255);not null;"`
	FunctionClaimID int64          `gorm:"column:function_claim_id;"`
	FunctionClaim   SFunctionClaim `gorm:"foreignKey:FunctionClaimID;references:id;constraint:OnDelete:CASCADE"`
	CreatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
