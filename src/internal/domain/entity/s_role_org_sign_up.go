package entity

import (
	"time"

	"github.com/google/uuid"
)

type SRoleOrgSignUp struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	RoleName  string    `json:"role_name" gorm:"type:varchar(255);not null"`
	OrgCode   string    `json:"org_code" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
