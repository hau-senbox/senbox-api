package request

import (
	"errors"

	"github.com/google/uuid"
)

type UpdateStudentRequest struct {
	StudentID   string `json:"student_id" binding:"required"`
	StudentName string `json:"student_name" binding:"required"`
}

func (r *UpdateStudentRequest) Validate() error {
	if _, err := uuid.Parse(r.StudentID); err != nil {
		return errors.New("student_id must be a valid UUID")
	}
	return nil
}
