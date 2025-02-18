package request

type CreateUserConfigRequest struct {
	TopButtonConfig      string `json:"top_button_config" binding:"required"`
	StudentOutputSheetId string `json:"student_output_sheet_id" binding:"required"`
	TeacherOutputSheetId string `json:"teacher_output_sheet_id" binding:"required"`
}
