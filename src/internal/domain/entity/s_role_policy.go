package entity

import (
	"time"
)

type SRolePolicy struct {
	ID          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	PolicyName  string    `gorm:"type:varchar(255);not null;"`
	Description string    `gorm:"type:varchar(255);not null;default:''"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`

	Roles  []SRole      `gorm:"many2many:s_role_policy_roles;foreignKey:id;joinForeignKey:policy_id;references:id;joinReferences:role_id;"`
	Claims []SRoleClaim `gorm:"many2many:s_role_policy_claims;foreignKey:id;joinForeignKey:policy_id;references:id;joinReferences:claim_id;"`
}
