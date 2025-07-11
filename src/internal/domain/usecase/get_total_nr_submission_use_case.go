package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"gorm.io/gorm"
)

type GetTotalNrSubmissionByConditionUseCase struct {
	submissionRepository *repository.SubmissionRepository
}

func NewGetTotalNrSubmissionByConditionUseCase(db *gorm.DB) *GetTotalNrSubmissionByConditionUseCase {
	return &GetTotalNrSubmissionByConditionUseCase{
		submissionRepository: &repository.SubmissionRepository{
			DBConn: db,
		},
	}
}

type GetTotalNrSubmissionByConditionInput struct {
<<<<<<< HEAD
	FormID   uint64         `json:"form_id"`
	UserID   string         `json:"user_id"`
	Key      *string        `json:"key"`
	DB       *string        `json:"db"`
	TimeSort value.TimeSort `json:"time_sort"`
	Duration *value.TimeRange
=======
	FormID       uint64         `json:"form_id"`
	UserID       string         `json:"user_id"`
	Key          *string        `json:"key"`
	DB           *string        `json:"db"`
	TimeSort     value.TimeSort `json:"time_sort"`
	DateDuration *value.TimeRange
>>>>>>> develop
}

func (uc *GetTotalNrSubmissionByConditionUseCase) Execute(input GetTotalNrSubmissionByConditionInput) (*response.GetSubmissionTotalNrResponse, error) {
	param := repository.GetSubmissionByConditionParam{
<<<<<<< HEAD
		FormID:   input.FormID,
		UserID:   input.UserID,
		Key:      input.Key,
		DB:       input.DB,
		TimeSort: input.TimeSort,
		Duration: input.Duration,
=======
		FormID:       input.FormID,
		UserID:       input.UserID,
		Key:          input.Key,
		DB:           input.DB,
		TimeSort:     input.TimeSort,
		DateDuration: input.DateDuration,
>>>>>>> develop
	}

	return uc.submissionRepository.GetTotalNrSubmissionByCondition(param)
}
