package entity

import "time"

type SDeviceFormDataset struct {
	ID                       string    `gorm:"type:varchar(255);primary_key;not null"`
	DeviceId                 string    `gorm:"type:varchar(36);primary_key;not null"`
	Set                      string    `gorm:"type:mediumtext;not null;"`
	QuestionDate             string    `gorm:"type:mediumtext;default:''"`
	QuestionTime             string    `gorm:"type:mediumtext;default:''"`
	QuestionDateTime         string    `gorm:"type:mediumtext;default:''"`
	QuestionDurationForward  string    `gorm:"type:mediumtext;default:''"`
	QuestionDurationBackward string    `gorm:"type:mediumtext;default:''"`
	QuestionScale            string    `gorm:"type:mediumtext;default:''"`
	QuestionQRCode           string    `gorm:"type:mediumtext;default:''"`
	QuestionSelection        string    `gorm:"type:mediumtext;default:''"`
	QuestionText             string    `gorm:"type:mediumtext;default:''"`
	QuestionCount            string    `gorm:"type:mediumtext;default:''"`
	QuestionNumber           string    `gorm:"type:mediumtext;default:''"`
	QuestionPhoto            string    `gorm:"type:mediumtext;default:''"`
	QuestionMultipleChoice   string    `gorm:"type:mediumtext;default:''"`
	QuestionButtonCount      string    `gorm:"type:mediumtext;default:''"`
	QuestionSingleChoice     string    `gorm:"type:mediumtext;default:''"`
	QuestionButtonList       string    `gorm:"type:mediumtext;default:''"`
	QuestionMessageBox       string    `gorm:"type:mediumtext;default:''"`
	QuestionShowPic          string    `gorm:"type:mediumtext;default:''"`
	QuestionButton           string    `gorm:"type:mediumtext;default:''"`
	QuestionPlayVideo        string    `gorm:"type:mediumtext;default:''"`
	QuestionQRCodeFront      string    `gorm:"type:mediumtext;default:''"`
	QuestionChoiceToggle     string    `gorm:"type:mediumtext;default:''"`
	QuestionSignature        string    `gorm:"type:mediumtext;default:''"`
	QuestionWeb              string    `gorm:"type:mediumtext;default:''"`
	QuestionDraggableList    string    `gorm:"type:mediumtext;default:''"`
	QuestionSendMessage      string    `gorm:"type:mediumtext;default:''"`
	CreatedAt                time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt                time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
