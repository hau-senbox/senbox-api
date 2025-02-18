package request

type UpdateDeviceRequestScreenButton struct {
	ButtonType  string `json:"button_type" binding:"required" enums:"scan,list"`
	ButtonValue string `json:"button_value" binding:"required"`
}

type UpdateDeviceRequest struct {
	Name              *string `json:"name"`
	DeactivateMessage *string `json:"deactivate_message"`
	ButtonType        *string `json:"button_type" enums:"scan,list"`
}

type UpdateDeviceRequestUserInfoV2 struct {
	UserInfo1Old    string `json:"user_info_1_old" binding:"required"`
	UserInfo2Old    string `json:"user_info_2_old" binding:"required"`
	UserInfo3Old    string `json:"user_info_3_old" binding:"required"`
	UserInfo1       string `json:"user_info_1" binding:"required"`
	UserInfo1ID     string `json:"user_info_1_id" binding:"required"`
	UserInfo1Prefix string `json:"user_info_1_prefix" binding:"required"`
	UserInfo2       string `json:"user_info_2" binding:"required"`
	UserInfo2ID     string `json:"user_info_2_id" binding:"required"`
	UserInfo2Prefix string `json:"user_info_2_prefix" binding:"required"`
	UserInfo3       string `json:"user_info_3" binding:"required"`
	UserInfo3ID     string `json:"user_info_3_id" binding:"required"`
	UserInfo3Prefix string `json:"user_info_3_prefix" binding:"required"`
}

type UpdateDeviceRequestV2 struct {
	Name                 *string                          `json:"name"`
	Note                 *string                          `json:"note"`
	Status               *string                          `json:"status"`
	OutputSpreadsheetUrl *string                          `json:"output_spreadsheet_url"`
	ButtonUrl            *string                          `json:"button_url"`
	Message              *string                          `json:"message"`
	UserInfo             *UpdateDeviceRequestUserInfoV2   `json:"user_info"`
	ScreenButton         *UpdateDeviceRequestScreenButton `json:"screen_button"`
}
