package response

type GetMenuSectionResponse struct {
	SectionName string              `json:"section_name"`
	SectionID   string              `json:"section_id"`
	Components  []ComponentResponse `json:"components"`
}
