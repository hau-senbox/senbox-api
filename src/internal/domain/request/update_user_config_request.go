package request

type UpdateUserConfigRequest struct {
	ID                   uint   `json:"id" binding:"required"`
	TopButtonConfig      string `json:"top_button_config" binding:"required"`
	ProfilePictureUrl    string `json:"profile_picture_url" binding:"required"`
	StudentOutputSheetId string `json:"student_output_sheet_id" binding:"required"`
	TeacherOutputSheetId string `json:"teacher_output_sheet_id" binding:"required"`
}
