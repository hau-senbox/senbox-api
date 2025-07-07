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
	FormID      uint64         `json:"form_id"`
	UserID      string         `json:"user_id"`
	QuestionKey string         `json:"question_key"`
	QuestionDB  string         `json:"question_db"`
	TimeSort    value.TimeSort `json:"time_sort"`
}

func (uc *GetTotalNrSubmissionByConditionUseCase) Execute(input GetTotalNrSubmissionByConditionInput) (*response.GetSubmissionTotalNrResponse, error) {
	param := repository.GetSubmissionByConditionParam{
		FormID:      input.FormID,
		UserID:      input.UserID,
		QuestionKey: input.QuestionKey,
		QuestionDB:  input.QuestionDB,
		TimeSort:    input.TimeSort,
	}

	return uc.submissionRepository.GetTotalNrSubmissionByCondition(param)
}
