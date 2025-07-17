package request

type UploadSectionMenuRequest []SectionMenuItem

type SectionMenuItem struct {
	SectionID  string                       `json:"section_id"`
	Components []CreateMenuComponentRequest `json:"components"`
}
