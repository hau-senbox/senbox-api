package usecase

import (
	"sen-global-api/internal/data/repository"
)

type GetUserTokenFCMUseCase struct {
	*repository.UserTokenFCMRepository
}

func (receiver *GetUserTokenFCMUseCase) GetAllToken(userID string) ([]string, error) {

	users, err := receiver.UserTokenFCMRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var tokens []string
	for _, user := range users {
		tokens = append(tokens, user.FCMToken)
	}

	return tokens, nil
	
}
