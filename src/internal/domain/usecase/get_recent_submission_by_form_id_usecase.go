package usecase

import (
	"encoding/json"
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/value"

	"gorm.io/gorm"
)

type GetRecentSubmissionByFormIDUseCase struct {
	formRepository       *repository.FormRepository
	submissionRepository *repository.SubmissionRepository
}

func NewGetRecentSubmissionByFormIDUseCase(db *gorm.DB) *GetRecentSubmissionByFormIDUseCase {
	return &GetRecentSubmissionByFormIDUseCase{
		formRepository: &repository.FormRepository{
			DBConn:                 db,
			DefaultRequestPageSize: 0,
		},
		submissionRepository: &repository.SubmissionRepository{
			DBConn: db,
		},
	}
}

type RecentSubmission struct {
	Items []RecentSubmissionItem `json:"items" binding:"required"`
}

type RecentSubmissionItem struct {
	QuestionID string `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

func (g *GetRecentSubmissionByFormIDUseCase) Execute(formID string, userID string) ([]RecentSubmissionItem, error) {
	form, err := g.formRepository.GetFormByQRCode(formID)
	if err != nil {
		return []RecentSubmissionItem{}, err
	}

	if form.Type != value.FormType_SelfRemember {
		return []RecentSubmissionItem{}, errors.New("this form does not support recent submission")
	}

	s, err := g.submissionRepository.FindRecentByFormID(form.ID, userID)
	if err != nil {
		return []RecentSubmissionItem{}, nil
	}

	var result RecentSubmission
	err = json.Unmarshal(s.SubmissionData, &result)
	if err != nil {
		return []RecentSubmissionItem{}, err
	}

	return result.Items, nil
}
