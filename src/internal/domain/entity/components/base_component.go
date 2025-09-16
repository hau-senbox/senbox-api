package components

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type IComponent interface {
	GetID() uuid.UUID
	GetName() string
	SetName(name string)
	GetType() ComponentType
	GetKey() string
	SetKey(key string)
	GetValue() datatypes.JSON
	SetValue(value datatypes.JSON)
	GetSectionID() string
	SetSectionID(SectionID string)

	GetComponent() *Component
	NormalizeValue() error
	JSONValue() (datatypes.JSON, error)
}

type Component struct {
	IComponent `gorm:"-" json:"-"`
	ID         uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	Name       string         `gorm:"type:varchar(255);not null" json:"name"`
	Type       ComponentType  `gorm:"type:varchar(255);not null" json:"type"`
	Key        string         `gorm:"type:varchar(255);not null;default:''" json:"key"`
	Value      datatypes.JSON `gorm:"type:json;not null;default:'{}'" json:"value"`
	SectionID  string         `gorm:"type:varchar(255);not null;default:''" json:"section_id"`
	Language   uint           `gorm:"type:int;not null;default:1" json:"language"`
}

func (component *Component) GetID() uuid.UUID {
	return component.ID
}

func (component *Component) GetName() string {
	return component.Name
}

func (component *Component) SetName(name string) {
	component.Name = name
}

func (component *Component) GetType() ComponentType {
	return component.Type
}

func (component *Component) GetKey() string {
	return component.Key
}

func (component *Component) SetKey(key string) {
	component.Key = key
}

func (component *Component) GetValue() datatypes.JSON {
	return component.Value
}

func (component *Component) SetValue(value datatypes.JSON) {
	component.Value = value
}

func (component *Component) GetComponent() *Component {
	return component
}

func (component *Component) GetSectionID() string {
	return component.SectionID
}

func (component *Component) SetSectionID(SectionID string) {
	component.SectionID = SectionID
}

type ComponentInnerValue struct {
	Visible      bool   `json:"visible"`
	Icon         string `json:"icon"`
	Color        string `json:"color"`
	FormQR       string `json:"form_qr"`
	Url          string `json:"url"`
	ShowedTop    *bool  `json:"showed_top"`
	ShowedBottom *bool  `json:"showed_bottom"`
}

type ComponentFullValue struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Type         string              `json:"type"`
	Key          string              `json:"key"`
	Value        ComponentInnerValue `json:"value"`
	SectionID    string              `json:"section_id"`
	Icon         string              `json:"icon"`
	Visible      bool                `json:"visible"`
	Color        string              `json:"color"`
	FormQR       string              `json:"form_qr"`
	Url          string              `json:"url"`
	ShowedTop    bool                `json:"showed_top"`
	ShowedBottom bool                `json:"showed_bottom"`
}

func (c *ComponentInnerValue) NormalizeDefault() {
	defaultTrue := true

	if c.ShowedTop == nil {
		c.ShowedTop = &defaultTrue
	}
	if c.ShowedBottom == nil {
		c.ShowedBottom = &defaultTrue
	}
}
