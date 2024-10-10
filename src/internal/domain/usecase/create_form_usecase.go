package usecase

import (
	log "github.com/sirupsen/logrus"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type CreateFormUseCase struct {
	*repository.FormRepository
	*repository.FormQuestionRepository
}

func (receiver *CreateFormUseCase) CreateForm(req request.CreateFormRequest) (*entity.SForm, error) {
	form, err := receiver.FormRepository.Create(&entity.SForm{
		Note: req.FormName,
	})
	if err != nil {
		return nil, err
	}

	formQuestions, err := receiver.FormQuestionRepository.CreateFormQuestions(form.FormId, req.Questions)
	if err != nil {
		return nil, err
	}

	log.Debug(formQuestions)

	return form, nil
}
