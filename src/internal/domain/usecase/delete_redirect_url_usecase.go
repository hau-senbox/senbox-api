package usecase

import "sen-global-api/internal/data/repository"

type DeleteRedirectUrlUseCase struct {
	*repository.RedirectUrlRepository
}

func (receiver *DeleteRedirectUrlUseCase) Delete(id uint64) error {
	return receiver.RedirectUrlRepository.Delete(id)
}
