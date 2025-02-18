package migrations

import (
	"encoding/json"
	"errors"
	"sen-global-api/internal/domain/value"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SLegacyDevice struct {
	DeviceId             string             `gorm:"type:varchar(36);primary_key;not null"`
	DeviceName           string             `gorm:"type:varchar(255);not null;default:''"`
	Attributes           datatypes.JSON     `gorm:"type:json;not null;default:'{}'"`
	PrimaryUserInfo      string             `gorm:"column:primary_user_info;type:varchar(255);not null;"`
	SecondaryUserInfo    string             `gorm:"column:secondary_user_info;type:varchar(255);not null"`
	ScreenButton         datatypes.JSON     `gorm:"type:json;not null;default:'{}'"`
	Status               value.DeviceStatus `gorm:"type:tinyint(1);not null;default:1"`
	ProfilePictureUrl    string             `gorm:"type:varchar(255);"`
	SpreadsheetId        string             `gorm:"type:varchar(255);not null;"`
	TeacherSpreadsheetId string             `gorm:"type:varchar(255);not null;default:''"`
	Message              string             `gorm:"type:varchar(255);not null;default:''"`
	ButtonUrl            string             `gorm:"type:varchar(255);not null;default:''"`
	Note                 string             `gorm:"type:varchar(255);not null;default:''"`
	AppVersion           string             `gorm:"type:varchar(255);not null;default:''"`
	RowNo                int                `gorm:"type:int;not null;default:0"`
	CreatedAt            time.Time          `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt            time.Time          `gorm:"default:CURRENT_TIMESTAMP;not null"`
}

type SDevice struct {
	DeviceId          string                 `gorm:"type:varchar(36);primary_key;not null"`
	DeviceName        string                 `gorm:"type:varchar(255);not null;default:''"`
	InputMode         value.InfoInputType    `gorm:"type:varchar(32);not null;default:1"`
	ScreenButtonType  value.ScreenButtonType `gorm:"type:varchar(16);not null;default:'scan'"`
	Status            value.DeviceMode       `gorm:"type:varchar(32);not null;default:1"`
	DeactivateMessage string                 `gorm:"type:varchar(255);not null;default:''"`
	ButtonUrl         string                 `gorm:"type:varchar(255);not null;default:''"`
	Note              string                 `gorm:"type:varchar(255);not null;default:''"`
	AppVersion        string                 `gorm:"type:varchar(255);not null;default:''"`
	RowNo             int                    `gorm:"type:int;not null;default:0"`
	CreatedAt         time.Time              `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt         time.Time              `gorm:"default:CURRENT_TIMESTAMP;not null"`
}

func MigrateDevices(db *gorm.DB) error {
	var devices []SLegacyDevice
	err := db.Table("s_device").Find(&devices).Error
	if err != nil {
		log.Info("s_device already migrated")
		return db.AutoMigrate(&SDevice{})
	}

	//err = db.AutoMigrate(&SDevice{})
	//if err != nil {
	//	return err
	//}

	var newDevices []SDevice
	for _, device := range devices {
		type DeviceAtt struct {
			InputMode string `json:"input_mode"`
		}

		var attributes DeviceAtt
		err = json.Unmarshal(device.Attributes, &attributes)
		if err != nil {
			return errors.New("no placeholder row found")
		}

		input, err := value.GetUserInfoInputTypeFromString(attributes.InputMode)
		if err != nil {
			return errors.New("no placeholder row found")
		}

		inputMode := value.InfoInputTypeKeyboard
		switch input {
		case value.UserInfoInputTypeKeyboard:
			inputMode = value.InfoInputTypeKeyboard
		case value.UserInfoInputTypeBarcode:
			inputMode = value.InfoInputTypeBarcode
		case value.UserInfoInputTypeBackOffice:
			inputMode = value.InfoInputTypeBackOffice
		}

		type ScreenButton struct {
			ButtonType  value.ButtonType `json:"button_type"`
			ButtonTitle string           `json:"button_title"`
		}

		var screenButtons ScreenButton
		buttonType := value.ScreenButtonType_Scan
		err = json.Unmarshal(device.ScreenButton, &screenButtons)
		if err != nil {
			return errors.New("no placeholder row found")
		}

		if screenButtons.ButtonType == value.ButtonTypeList {
			buttonType = value.ScreenButtonType_List
		}
		mode := value.DeviceModeSuspended
		switch device.Status {
		case value.DeviceStatus_Suspend:
			mode = value.DeviceModeSuspended
		case value.DeviceStatus_ModeT:
			mode = value.DeviceModeT
		case value.DeviceStatus_ModeP:
			mode = value.DeviceModeP
		case value.DeviceStatus_ModeS:
			mode = value.DeviceModeS
		case value.DeviceStatus_Deactive:
			mode = value.DeviceModeDeactivated
		case value.DeviceStatus_ModeL:
			mode = value.DeviceModeL
		}
		var newDevice SDevice
		newDevice.DeviceId = device.DeviceId
		newDevice.DeviceName = device.DeviceName
		newDevice.InputMode = inputMode
		newDevice.ScreenButtonType = buttonType
		newDevice.Status = mode
		newDevice.DeactivateMessage = device.Message
		newDevice.ButtonUrl = device.ButtonUrl
		newDevice.Note = device.Note
		newDevice.AppVersion = device.AppVersion
		newDevice.RowNo = device.RowNo
		newDevice.CreatedAt = device.CreatedAt

		newDevices = append(newDevices, newDevice)
	}

	err = db.Exec("DROP TABLE s_device").Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&SDevice{})
	if err != nil {
		return err
	}

	//check empty
	if len(newDevices) == 0 {
		return nil
	}

	return db.Save(&newDevices).Error
}
