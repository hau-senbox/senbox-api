package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SRoleOrgSignUp struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key"`
	RoleName  string    `json:"role_name" gorm:"type:varchar(255);not null"`
	OrgCode   string    `json:"org_code" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (r *SRoleOrgSignUp) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}
