package components

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

var _ IComponent = (*ButtonFormComponent)(nil)

type ButtonFormComponent struct {
	*Component
	Icon    string `json:"icon"`
	Visible bool   `json:"visible"`
	Color   string `json:"color"`
	FormQR  string `json:"form_qr"`
	//Required bool   `json:"required"`
}

func NewButtonFormComponent() *ButtonFormComponent {
	base := &Component{}
	button := &ButtonFormComponent{}
	button.Component = base
	button.ID = uuid.New()
	button.Type = ButtonForm

	base.IComponent = button

	return button
}

func (c *ButtonFormComponent) NormalizeValue() (err error) {
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
	formQR := dataMap["form_qr"].(string)

	c.Icon = icon
	c.Visible = visible
	c.Color = color
	c.FormQR = formQR

	value, err := c.JSONValue()
	if err != nil {
		return err
	}

	c.Value = value

	return nil
}

func (c *ButtonFormComponent) JSONValue() (datatypes.JSON, error) {
	j, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return j, nil
}
