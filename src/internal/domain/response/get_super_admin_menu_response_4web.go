package response

type GetSuperAdminMenuResponse4Web struct {
	Top    []GetMenus4Web `json:"top"`
	Bottom []GetMenus4Web `json:"bottom"`
}
