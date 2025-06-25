package repository

import (
	"errors"
	"github.com/google/uuid"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/parameters"
	"sen-global-api/internal/domain/request"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FormQuestionRepository struct {
	DBConn *gorm.DB
}

func (receiver *FormQuestionRepository) CreateFormQuestions(formID uint64, questionItems []request.CreateFormQuestionItem) ([]entity.SFormQuestion, error) {
	var formQuestions []entity.SFormQuestion
	for _, questionItem := range questionItems {
		formQuestions = append(formQuestions, entity.SFormQuestion{
			FormID:         formID,
			QuestionID:     uuid.MustParse(questionItem.QuestionID),
			Order:          questionItem.Order,
			AnswerRequired: questionItem.AnswerRequired,
			AnswerRemember: questionItem.AnswerRemember,
		})
	}

	err := receiver.DBConn.Table("s_form_question").Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "form_id"}, {Name: "question_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"order", "answer_required", "answer_remember"}),
		}).Create(&formQuestions).Error
	if err != nil {
		return nil, err
	}

	return formQuestions, nil
}

func (receiver *FormQuestionRepository) Update(form *entity.SForm, questions []entity.SQuestion, rawQuestions []parameters.RawQuestion) error {
	if form == nil {
		return errors.New("form is not found")
	}
	formQuestions := make([]entity.SFormQuestion, 0)
	for _, rawQuestion := range rawQuestions {
		formQuestions = append(formQuestions, entity.SFormQuestion{
			FormID:         form.ID,
			QuestionID:     uuid.MustParse(rawQuestion.QuestionID),
			Order:          rawQuestion.RowNumber,
			AnswerRequired: strings.ToLower(rawQuestion.AnswerRequired) == "true",
		})
	}

	err := receiver.DBConn.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec("SET FOREIGN_KEY_CHECKS=0;").Error
		if err != nil {
			return err
		}

		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "form_id"}, {Name: "question_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"order", "answer_required"}),
		}).Create(&formQuestions).Error
	})

	return err
}

func (receiver *FormQuestionRepository) GetFormQuestionsByForm(form entity.SForm) ([]entity.SFormQuestion, error) {
	var formQuestions []entity.SFormQuestion
	err := receiver.DBConn.Where("form_id = ?", form.ID).Find(&formQuestions).Error
	if err != nil {
		return nil, err
	}

	return formQuestions, nil
}

func (receiver *FormQuestionRepository) DeleteByFormID(formID uint64) error {
	return receiver.DBConn.Where("form_id = ?", formID).Delete(&entity.SFormQuestion{}).Error
}
