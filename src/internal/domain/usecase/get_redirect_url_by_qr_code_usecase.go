package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetRedirectUrlByQRCodeUseCase struct {
	*repository.RedirectUrlRepository
}

func (receiver *GetRedirectUrlByQRCodeUseCase) GetByQRCode(qrCode string) (*entity.SRedirectUrl, error) {
	return receiver.RedirectUrlRepository.GetByQRCode(qrCode)
}
