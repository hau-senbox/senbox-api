package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/value"

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
	FormID   uint64         `json:"form_id"`
	UserID   string         `json:"user_id"`
	Key      *string        `json:"key"`
	DB       *string        `json:"db"`
	TimeSort value.TimeSort `json:"time_sort"`
	Quantity *string        `json:"quantity"`
}

func (uc *GetSubmissionByConditionUseCase) Execute(input GetSubmissionByConditionInput) (*[]repository.SubmissionDataItem, error) {
	param := repository.GetSubmissionByConditionParam{
		FormID:   input.FormID,
		UserID:   input.UserID,
		Key:      input.Key,
		DB:       input.DB,
		TimeSort: value.TimeSort(input.TimeSort),
		Quantity: input.Quantity,
	}

	items, err := uc.submissionRepository.GetSubmissionByCondition(param)

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (uc *GetSubmissionByConditionUseCase) ExGetSubmission4ListRes(input GetSubmissionByConditionInput) (*[]repository.SubmissionDataItem, error) {
	param := repository.GetSubmissionByConditionParam{
		FormID:   input.FormID,
		UserID:   input.UserID,
		Key:      input.Key,
		DB:       input.DB,
		TimeSort: value.TimeSort(input.TimeSort),
		Quantity: input.Quantity,
	}

	items, err := uc.submissionRepository.GetSubmissionByCondition(param)

	if err != nil {
		return nil, err
	}

	return items, nil
}
