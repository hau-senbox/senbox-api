package usecase

import (
	"encoding/json"
	"log"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"strings"

	"gorm.io/gorm"
)

type GetSubmissionChildProfileUseCase struct {
	formRepository       *repository.FormRepository
	submissionRepository *repository.SubmissionRepository
	questionRepository   *repository.QuestionRepository
}

func NewGetSubmissionChildProfileUseCase(db *gorm.DB) *GetSubmissionChildProfileUseCase {
	return &GetSubmissionChildProfileUseCase{
		formRepository: &repository.FormRepository{
			DBConn:                 db,
			DefaultRequestPageSize: 0,
		},
		questionRepository: &repository.QuestionRepository{
			DBConn: db,
		},
		submissionRepository: &repository.SubmissionRepository{
			DBConn: db,
		},
	}
}

func (uc *GetSubmissionChildProfileUseCase) Execute(input repository.GetSubmissionChildProfileParam) ([]repository.SubmissionDataItem, error) {

	questions, err := uc.questionRepository.GetQuestionsByFormID(input.FormID)
	if err != nil {
		return nil, nil
	}
	var rawQuestions = make([]response.QuestionListData, 0)
	for _, question := range questions {
		var att response.QuestionAttributes
		if err := json.Unmarshal(question.Attributes, &att); err != nil {
			log.Printf("Invalid attributes for question %s: %v", question.ID, err)
			continue
		}

		// filter Question
		if question.QuestionType == value.GetStringValue(value.QuestionSendNotification) {
			continue
		}

		if question.QuestionKey == "" && question.QuestionDB == "" {
			continue
		}
		q := response.QuestionListData{
			QuestionID:     question.ID,
			QuestionType:   question.QuestionType,
			Question:       question.Question,
			Attributes:     att,
			Order:          question.Order,
			AnswerRequired: question.AnswerRequired,
			AnswerRemember: question.AnswerRemember,
			Enabled:        question.EnableOnMobile == value.QuestionForMobile_Enabled,
			QuestionKey:    question.QuestionKey,
			QuestionDB:     question.QuestionDB,
		}

		rawQuestions = append(rawQuestions, q)
	}

	items, err := uc.submissionRepository.GetSubmissionChildProfile(input)
	if err != nil {
		return nil, err
	}

	// Tạo map để tra cứu item theo QuestionID
	existingMap := make(map[string]bool)
	for _, item := range items {
		existingMap[item.QuestionKey] = true
	}

	// Thêm các câu hỏi chưa có trong submission (tức là chưa trả lời)
	for _, q := range rawQuestions {
		if !existingMap[q.QuestionKey] {
			items = append(items, repository.SubmissionDataItem{
				SubmissionID: 0,
				QuestionID:   q.QuestionID,
				QuestionKey:  q.QuestionKey,
				QuestionDB:   q.QuestionDB,
				Question:     q.Question,
				Answer:       "",
			})
		}
	}

	// loc key, db == ""
	filteredItems := make([]repository.SubmissionDataItem, 0, len(items))
	for _, item := range items {
		if strings.TrimSpace(item.QuestionKey) == "" && strings.TrimSpace(item.QuestionDB) == "" {
			continue
		}
		filteredItems = append(filteredItems, item)
	}

	return filteredItems, nil
}
