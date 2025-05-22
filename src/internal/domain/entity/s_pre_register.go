package entity

type SPreRegister struct {
	Email string `gorm:"column:email;primaryKey"`
}
