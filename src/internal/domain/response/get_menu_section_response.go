package response

type GetMenuSectionResponse struct {
	SectionName string              `json:"section_name"`
	SectionID   string              `json:"section_id"`
	MenuIconKey string              `json:"menu_icon_key"`
	Components  []ComponentResponse `json:"components"`
}
