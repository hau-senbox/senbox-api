package response

import "sen-global-api/internal/domain/value"

type DeviceComponentValuesResponse struct {
	ID           int             `json:"id"`
	Setting      SettingResponse `json:"setting"`
	Organization string          `json:"organization_id"`
}

type ComponentPositionValues string

type ComponentPositionAttribute struct {
	Name  string                  `json:"name"`
	Value ComponentPositionValues `json:"value"`
}

type SettingResponse struct {
	ComponentPositionAttribute []ComponentPositionAttribute `json:"component_attributes"`
	ComponentSettings          []ComponentSettingResponse   `json:"component_settings"`
}

type ComponentSettingResponse struct {
	ComponentName         string                        `json:"component_name"`
	ComponentType         string                        `json:"component_type"`
	ComponentPosition     *ComponentPositionValues      `json:"component_position"`
	ComponentPositionRoot ComponentPositionValues       `json:"component_position_root"`
	ComponentConditions   *[]ComponentConditionResponse `json:"component_condition"`
}

type ComponentConditionResponse struct {
	ConditionName string                   `json:"condition_name"`
	ConditionKey  value.DeviceConditionKey `json:"condition_key"`
	DelayDisplay  *int                     `json:"delay_display"`
	Priority      int                      `json:"priority"`
}
