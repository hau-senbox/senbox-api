package entity

import "github.com/google/uuid"

type SUserOrg struct {
	UserID         uuid.UUID     `gorm:"column:user_id;primary_key"`
	User           SUserEntity   `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE;"`
	OrganizationID uuid.UUID     `gorm:"column:organization_id;primary_key"`
	Organization   SOrganization `gorm:"foreignKey:OrganizationID;references:id;constraint:OnDelete:CASCADE;"`
	UserNickName   string        `gorm:"type:varchar(255);not null;default:''"`
	IsManager      bool          `gorm:"type:tinyint(1);not null;default:0"`
}
