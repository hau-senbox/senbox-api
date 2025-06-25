package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetFormByIDUseCase struct {
	*repository.FormRepository
}

func (receiver *GetFormByIDUseCase) GetFormByID(formID int) (*entity.SForm, error) {
	return receiver.FormRepository.GetFormByID(uint64(formID))
}

func (receiver *GetFormByIDUseCase) GetFormByQRCode(qrCode string) (*entity.SForm, error) {
	return receiver.FormRepository.GetFormByQRCode(qrCode)
}
