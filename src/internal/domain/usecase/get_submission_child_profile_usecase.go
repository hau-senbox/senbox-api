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

type GetSubmission4MemoriesFormUseCase struct {
	formRepository       *repository.FormRepository
	submissionRepository *repository.SubmissionRepository
	questionRepository   *repository.QuestionRepository
}

func NewGetSubmission4MemoriesFormUseCase(db *gorm.DB) *GetSubmission4MemoriesFormUseCase {
	return &GetSubmission4MemoriesFormUseCase{
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

func (uc *GetSubmission4MemoriesFormUseCase) Execute(input repository.GetSubmission4MemoriesFormParam) ([]repository.SubmissionDataItem, error) {

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

		if question.Key == "" && question.DB == "" {
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
			Key:            question.Key,
			DB:             question.DB,
		}

		rawQuestions = append(rawQuestions, q)
	}

	items, err := uc.submissionRepository.GetSubmission4MemoriesForm(input)
	if err != nil {
		return nil, err
	}

	// Tạo map để tra cứu item theo QuestionID
	existingMap := make(map[string]bool)
	for _, item := range items {
		existingMap[item.DB] = true
	}

	// Thêm các câu hỏi chưa có trong submission (tức là chưa trả lời)
	for _, q := range rawQuestions {
		if !existingMap[q.DB] {
			items = append(items, repository.SubmissionDataItem{
				SubmissionID: 0,
				QuestionID:   q.QuestionID,
				Key:          q.Key,
				DB:           q.DB,
				Question:     q.Question,
				Answer:       "",
			})
		}
	}

	// loc key, db == ""
	filteredItems := make([]repository.SubmissionDataItem, 0, len(items))
	for _, item := range items {
		if strings.TrimSpace(item.Key) == "" && strings.TrimSpace(item.DB) == "" {
			continue
		}
		filteredItems = append(filteredItems, item)
	}

	for i := range filteredItems {
		filteredItems[i].QuestionData = response.QuestionListData{
			QuestionType: "memory_text",
			Question:     filteredItems[i].Question,
			Attributes:   response.QuestionAttributes{},
			Order:        i,
			Key:          filteredItems[i].Key,
			DB:           filteredItems[i].DB,
		}
	}

	return filteredItems, nil
}
