package response

type ParentResponseBase struct {
	ParentID   string              `json:"id"`
	ParentName string              `json:"name"`
	Avatar     string              `json:"avatar,omitempty"`
	AvatarURL  string              `json:"avatar_url,omitempty"`
	Menus      []ComponentResponse `json:"components,omitempty"`
	CustomID   string              `json:"custom_id"`
}
