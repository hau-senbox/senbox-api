package usecase

import (
	"sen-global-api/internal/data/repository"

	"gorm.io/gorm"
)

type GetSubmissionByConditionUseCase struct {
	formRepository       *repository.FormRepository
	submissionRepository *repository.SubmissionRepository
}

func NewGetSubmissionByConditionUseCase(db *gorm.DB) *GetSubmissionByConditionUseCase {
	return &GetSubmissionByConditionUseCase{
		formRepository: &repository.FormRepository{
			DBConn:                 db,
			DefaultRequestPageSize: 0,
		},
		submissionRepository: &repository.SubmissionRepository{
			DBConn: db,
		},
	}
}

type GetSubmissionByConditionInput struct {
	FormID      uint64
	UserID      string
	QuestionKey string
	QuestionDB  string
}

func (uc *GetSubmissionByConditionUseCase) Execute(input GetSubmissionByConditionInput) ([]repository.SubmissionDataItem, error) {
	param := repository.GetSubmissionByConditionParam{
		FormID:      input.FormID,
		UserID:      input.UserID,
		QuestionKey: input.QuestionKey,
		QuestionDB:  input.QuestionDB,
	}

	items, err := uc.submissionRepository.GetSubmissionByCondition(param)
	if err != nil {
		return nil, err
	}

	return items, nil
}
