package entity

type OrganizationMenuTemplate struct {
	ID             string `json:"id" gorm:"type:char(36);primary_key"`
	OrganizationID string `json:"organization_id" gorm:"type:char(36);not null"`
	ComponentID    string `json:"component_id" gorm:"type:char(36);not null"`
	SectionID      string `json:"section_id" gorm:"type:char(36);not null"`
}
