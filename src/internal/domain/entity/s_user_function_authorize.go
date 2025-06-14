package entity

import "github.com/google/uuid"

type SUserFunctionAuthorize struct {
	UserID                    uuid.UUID                `gorm:"column:user_id;primary_key"`
	User                      SUserEntity              `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE;"`
	FunctionClaimId           int64                    `gorm:"column:function_claim_id;primary_key"`
	FunctionClaim             SFunctionClaim           `gorm:"foreignKey:FunctionClaimId;references:id;constraint:OnDelete:CASCADE"`
	FunctionClaimPermissionId int64                    `gorm:"column:function_claim_permission_id;"`
	FunctionClaimPermission   SFunctionClaimPermission `gorm:"foreignKey:FunctionClaimPermissionId;references:id;constraint:OnDelete:CASCADE"`
}
