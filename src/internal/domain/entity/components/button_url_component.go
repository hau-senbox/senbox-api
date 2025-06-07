package components

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

var _ IComponent = (*ButtonURLComponent)(nil)

type ButtonURLComponent struct {
	*Component
	Icon    string `json:"icon"`
	Visible bool   `json:"visible"`
	Color   string `json:"color"`
	URL     string `json:"url"`
	//Required bool   `json:"required"`
}

func NewButtonURLComponent() *ButtonURLComponent {
	base := &Component{}
	button := &ButtonURLComponent{}
	button.Component = base
	button.ID = uuid.New()
	button.Type = ButtonURL

	base.IComponent = button

	return button
}

func (c *ButtonURLComponent) NormalizeValue() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error normalizing value: %v", r)
		}
	}()

	var data interface{}
	err = json.Unmarshal(c.Value, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal attributes: %w", err)
	}

	if data == nil {
		return nil
	}

	dataMap := data.(map[string]interface{})

	icon := dataMap["icon"].(string)
	visible := dataMap["visible"].(bool)
	color := dataMap["color"].(string)
	url := dataMap["url"].(string)

	c.Icon = icon
	c.Visible = visible
	c.Color = color
	c.URL = url

	value, err := c.JSONValue()
	if err != nil {
		return err
	}

	c.Value = value

	return nil
}

func (c *ButtonURLComponent) JSONValue() (datatypes.JSON, error) {
	j, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return j, nil
}
