package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type CreateChildForParentUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *CreateChildForParentUseCase) CreateChildForParent(parentID string, req request.CreateChildForParentRequest) error {
	return receiver.UserEntityRepository.CreateChildForParent(parentID, req)
}
