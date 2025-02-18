package entity

type SRolePolicyClaims struct {
	PolicyId  int64       `gorm:"column:policy_id;primary_key"`
	Policy    SRolePolicy `gorm:"foreignKey:PolicyId;references:id;constraint:OnDelete:CASCADE;"`
	ClaimId   int64       `gorm:"column:claim_id;primary_key"`
	RoleClaim SRoleClaim  `gorm:"foreignKey:ClaimId;references:id;constraint:OnDelete:CASCADE"`
}
