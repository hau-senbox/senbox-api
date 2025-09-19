package request

import (
	"fmt"

	"github.com/google/uuid"
)

type UploadSectionMenuStudentRequest StudentSectionMenuItem

type StudentSectionMenuItem struct {
	LanguageID         uint                         `json:"language_id" binding:"required"`
	StudentID          string                       `json:"student_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}

// Validate UUID format in DeleteComponentIDs
func (s *StudentSectionMenuItem) Validate() error {
	for _, id := range s.DeleteComponentIDs {
		if _, err := uuid.Parse(id); err != nil {
			return fmt.Errorf("invalid UUID in delete_component_ids: %s", id)
		}
	}
	return nil
}
