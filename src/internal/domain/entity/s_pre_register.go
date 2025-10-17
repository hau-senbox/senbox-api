package entity

import "time"

type SPreRegister struct {
	Email      string     `gorm:"column:email;primaryKey"`
	DeviceID   string     `gorm:"column:device_id;type:varchar(100);default:''"`
	DeviceName string     `gorm:"column:device_name;type:varchar(100);default:''"`
	FormQR     string     `gorm:"column:form_qr;type:varchar(100);default:''"`
	CreatedAt  *time.Time `gorm:"column:created_at;default:null"`
}
