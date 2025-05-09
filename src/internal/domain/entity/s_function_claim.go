package entity

import (
	"time"
)

type SFunctionClaim struct {
	ID           int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	FunctionName string    `gorm:"type:varchar(255);not null;"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`

	ClaimPermissions []SFunctionClaimPermission `gorm:"foreignKey:function_claim_id;references:id;constraint:OnDelete:CASCADE"`
}
