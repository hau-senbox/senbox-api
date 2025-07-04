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

// type TimeShort string

// const (
// 	TimeShortLatest TimeShort = "latest"
// 	TimeShortOldest TimeShort = "oldest"
// )

type GetSubmissionByConditionInput struct {
	FormID      uint64              `json:"form_id"`
	UserID      string              `json:"user_id"`
	QuestionKey string              `json:"question_key"`
	QuestionDB  string              `json:"question_db"`
	TimeSort    repository.TimeSort `json:"time_sort"` // "latest" (default) or "oldest"
}

func (uc *GetSubmissionByConditionUseCase) Execute(input GetSubmissionByConditionInput) (repository.SubmissionDataItem, error) {
	param := repository.GetSubmissionByConditionParam{
		FormID:      input.FormID,
		UserID:      input.UserID,
		QuestionKey: input.QuestionKey,
		QuestionDB:  input.QuestionDB,
		TimeSort:    repository.TimeSort(input.TimeSort),
	}

	item, err := uc.submissionRepository.GetSubmissionByCondition(param)
	if err != nil {
		return repository.SubmissionDataItem{}, err
	}

	return item, nil
}
