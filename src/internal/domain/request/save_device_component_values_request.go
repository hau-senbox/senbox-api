package request

type SaveDeviceComponentValuesByOrganizationRequest struct {
	ID       *uint              `json:"id"`
	Settings SaveSettingRequest `json:"setting" binding:"required"`
	Organization  uint               `json:"organization_id" binding:"required"`
}

type SaveDeviceComponentValuesByDeviceRequest struct {
	ID       *uint              `json:"id"`
	Settings SaveSettingRequest `json:"setting" binding:"required"`
}

type ComponentPositionValues string

type ComponentPositionAttribute struct {
	Name  string                  `json:"name" binding:"required"`
	Value ComponentPositionValues `json:"value" binding:"required"`
}

type SaveSettingRequest struct {
	ComponentPositionAttributes []ComponentPositionAttribute  `json:"component_attributes" binding:"required"`
	ComponentSettings           []SaveComponentSettingRequest `json:"component_settings" binding:"required"`
}

type SaveComponentSettingRequest struct {
	ComponentName         string                            `json:"component_name" binding:"required"`
	ComponentType         string                            `json:"component_type" binding:"required"`
	ComponentPosition     *ComponentPositionValues          `json:"component_position" binding:"required"`
	ComponentPositionRoot ComponentPositionValues           `json:"component_position_root" binding:"required"`
	ComponentConditions   *[]SaveComponentConditionsRequest `json:"component_condition"`
}

type SaveComponentConditionsRequest struct {
	ConditionName  string  `json:"condition_name" binding:"required"`
	ConditionValue *string `json:"condition_value"`
	DelayDisplay   *int    `json:"delay_display"`
	Priority       int     `json:"priority" binding:"required"`
}
