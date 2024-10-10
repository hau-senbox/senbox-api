package response

type GetSettingMessageItem struct {
	Description string `json:"description" binding:"required"`
	Message     string `json:"message" binding:"required"`
}

type GetSettingMessageResponseData struct {
	Messages []GetSettingMessageItem `json:"messages" binding:"required"`
	FontSize *int                    `json:"font_size" binding:"required"`
}

type GetSettingMessageResponse struct {
	Data GetSettingMessageResponseData `json:"data"`
}
