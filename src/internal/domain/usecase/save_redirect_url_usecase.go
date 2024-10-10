package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type SaveRedirectUrlUseCase struct {
	*repository.RedirectUrlRepository
}

func (receiver *SaveRedirectUrlUseCase) Save(req request.SaveRedirectUrlRequest) (*entity.SRedirectUrl, error) {
	var url = entity.SRedirectUrl{
		QRCode:    req.QRCode,
		TargetUrl: req.TargetUrl,
	}
	if req.Password != "" {
		url.Password = &req.Password
	}
	return receiver.RedirectUrlRepository.Save(url)
}
