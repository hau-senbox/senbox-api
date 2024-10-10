package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetFormByIdUseCase struct {
	*repository.FormRepository
}

func (receiver *GetFormByIdUseCase) GetFormById(formId int) (*entity.SForm, error) {
	return receiver.FormRepository.GetFormById(uint64(formId))
}

func (receiver *GetFormByIdUseCase) GetFormByQRCode(qrCode string) (*entity.SForm, error) {
	return receiver.FormRepository.GetFormByQRCode(qrCode)
}
