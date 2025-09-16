package request

import (
	"fmt"

	"github.com/google/uuid"
)

type UploadSectionMenuDeviceRequest DeviceSectionMenuItem

type DeviceSectionMenuItem struct {
	Language           uint                         `json:"language" binding:"required"`
	DeviceID           string                       `json:"device_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}

// Validate UUID format in DeleteComponentIDs
func (s *DeviceSectionMenuItem) Validate() error {
	for _, id := range s.DeleteComponentIDs {
		if _, err := uuid.Parse(id); err != nil {
			return fmt.Errorf("invalid UUID in delete_component_ids: %s", id)
		}
	}
	return nil
}
