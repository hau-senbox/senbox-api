package usecase

import "sen-global-api/internal/data/repository"

type DeleteFormUseCase struct {
	*repository.FormRepository
}

func (receiver *DeleteFormUseCase) DeleteForm(formId uint64) error {
	err := receiver.FormRepository.DeleteForm(formId)
	if err != nil {
		return err
	}

	return nil
}
