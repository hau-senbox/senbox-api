package response

import "gorm.io/datatypes"

type UserSettingResponse struct {
	Key   string         `json:"key"`
	Value datatypes.JSON `json:"value"`
}
