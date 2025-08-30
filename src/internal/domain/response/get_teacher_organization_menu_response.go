package response

type GetTeacherOrganizationMenuResponse struct {
	Section     string              `json:"section_name"`
	MenuIconKey string              `json:"menu_icon_key"`
	Components  []ComponentResponse `json:"components"`
}
