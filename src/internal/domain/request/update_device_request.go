package request

type UpdateDeviceRequestUserInfo struct {
	UserInfo1 string `json:"user_info_1" binding:"required"`
	UserInfo2 string `json:"user_info_2" binding:"required"`
	UserInfo3 string `json:"user_info_3" binding:"required"`
}

type UpdateDeviceRequestScreenButton struct {
	ButtonType  string `json:"button_type" binding:"required" enums:"scan,list"`
	ButtonValue string `json:"button_value" binding:"required"`
}

type UpdateDeviceRequest struct {
	Name         *string                          `json:"name"`
	ButtonUrl    *string                          `json:"button_url"`
	Message      *string                          `json:"message"`
	UserInfo     *UpdateDeviceRequestUserInfo     `json:"user_info"`
	ScreenButton *UpdateDeviceRequestScreenButton `json:"screen_button"`
}
