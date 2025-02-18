package response

type UserConfigResponse struct {
	ID                   int64  `json:"id"`
	TopButtonConfig      string `json:"top_button_config"`
	StudentOutputSheetId string `json:"student_output_sheet_id"`
	TeacherOutputSheetId string `json:"teacher_output_sheet_id"`
}
