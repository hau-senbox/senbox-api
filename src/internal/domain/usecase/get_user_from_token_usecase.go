package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetUserFromTokenUseCase struct {
	repository.UserRepository
	repository.SessionRepository
}

func (c *GetUserFromTokenUseCase) GetUserFromToken(tokenString string) (*entity.SUser, error) {
	userId, err := c.SessionRepository.ExtractUserIdFromToken(tokenString)
	if err != nil {
		return nil, err
	}

	return c.UserRepository.FindUserById(userId)
}
