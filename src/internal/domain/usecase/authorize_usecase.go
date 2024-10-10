package usecase

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
)

type AuthorizeUseCase struct {
	repository.UserRepository
	repository.SessionRepository
}

func (receiver AuthorizeUseCase) LoginInputDao(req request.LoginInputReq) (response.LoginResponseData, error) {
	user := receiver.UserRepository.FindUserByUsername(req.LoginId)
	if user == nil {
		log.Info("No user has username matches", req.LoginId)
		return response.LoginResponseData{}, errors.New("user not found")
	}

	err := receiver.SessionRepository.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return response.LoginResponseData{}, errors.New("invalid username or password")
	}

	token, err := receiver.SessionRepository.GenerateToken(*user)
	if err != nil {
		return response.LoginResponseData{}, errors.New("cannot generate token")
	}

	//authMiddleware := jwtauth.JwtMiddleware()
	//token := authMiddleware.TokenGen(user.UserId)
	return *token, nil
}
