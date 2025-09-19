package request

import (
	"fmt"

	"github.com/google/uuid"
)

type UploadSectionMenuChildRequest ChildSectionMenuItem

type ChildSectionMenuItem struct {
	LanguageID         uint                         `json:"language_id" binding:"required"`
	ChildID            string                       `json:"child_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}

// Validate UUID format in DeleteComponentIDs
func (s *ChildSectionMenuItem) Validate() error {
	for _, id := range s.DeleteComponentIDs {
		if _, err := uuid.Parse(id); err != nil {
			return fmt.Errorf("invalid UUID in delete_component_ids: %s", id)
		}
	}
	return nil
}
