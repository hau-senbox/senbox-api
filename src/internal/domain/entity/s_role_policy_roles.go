package entity

type SRolePolicyRoles struct {
	PolicyId   int64       `gorm:"column:policy_id;primary_key"`
	RolePolicy SRolePolicy `gorm:"foreignKey:PolicyId;references:id;constraint:OnDelete:CASCADE;"`
	RoleId     int64       `gorm:"column:role_id;primary_key"`
	Role       SRole       `gorm:"foreignKey:RoleId;references:id;constraint:OnDelete:CASCADE"`
}
