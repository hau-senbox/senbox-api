package entity

import "github.com/google/uuid"

type SUserFunctionAuthorize struct {
	UserID                    uuid.UUID                `gorm:"column:user_id;primary_key"`
	User                      SUserEntity              `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE;"`
	FunctionClaimID           int64                    `gorm:"column:function_claim_id;primary_key"`
	FunctionClaim             SFunctionClaim           `gorm:"foreignKey:FunctionClaimID;references:id;constraint:OnDelete:CASCADE"`
	FunctionClaimPermissionID int64                    `gorm:"column:function_claim_permission_id;"`
	FunctionClaimPermission   SFunctionClaimPermission `gorm:"foreignKey:FunctionClaimPermissionID;references:id;constraint:OnDelete:CASCADE"`
}
