package response

import "sen-global-api/internal/domain/entity"

// common
type GetCommonMenuResponse struct {
	ChildID    string              `json:"child_id,omitempty"`
	Components []ComponentResponse `json:"components"`
}

type ComponentResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Key        string `json:"key"`
	Value      string `json:"value"`
	Order      int    `json:"order"`
	IsShow     bool   `json:"is_show"`
	LanguageID uint   `json:"language_id"`
}

type ComponentCommonMenu struct {
	ChildID    string              `json:"child_id,omitempty"`
	Components []ComponentResponse `json:"components"`
}

// by user
type GetCommonMenuByUserResponse struct {
	Components []ComponentCommonMenuByUser `json:"components"`
}

type ComponentCommonMenuByUser struct {
	ChildID   string            `json:"child_id,omitempty"`
	Component ComponentResponse `json:"component"`
}

type GetMenus4Web struct {
	Language entity.LanguageSetting `json:"language"`
	Menus    []ComponentResponse    `json:"menus"`
}
