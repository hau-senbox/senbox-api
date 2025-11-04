package request

type UploadUserSettingRequest struct {
	OwnerID   string `json:"owner_id" binding:"required"`
	OwnerRole string `json:"owner_role" binding:"required"`
	Key       string `json:"key" binding:"required"`
	Value     any    `json:"value" binding:"required"`
}

type UploadUserIsFirstLoginRequest struct {
	UserID       string `json:"user_id"`
	IsFirstLogin bool   `json:"is_first_login" binding:"required"`
}

type UploadUserWelcomeReminderRequest struct {
	UserID       string `json:"user_id"`
	IsEnabled    bool   `json:"is_enabled" binding:"required"`
	TimeReminder string `json:"time_reminder"`
}
