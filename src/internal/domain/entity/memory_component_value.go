package entity

type MemoryComponentValue struct {
	ComponentName string `gorm:"column:component_name;primary_key"`
	Value         string `gorm:"column:value;type:text;not null;default:''"`
}
