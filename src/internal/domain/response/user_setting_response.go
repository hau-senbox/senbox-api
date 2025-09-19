package response

import "gorm.io/datatypes"

type UserSettingResponse struct {
	LimitDeviceLogin datatypes.JSON `json:"limit_device_login"`
	AppLanguage      datatypes.JSON `json:"app_language"`
}
