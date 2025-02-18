package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"

	log "github.com/sirupsen/logrus"
)

type AuthorizeUseCase struct {
	*repository.UserRepository
	*repository.UserEntityRepository
	repository.SessionRepository
}

func (receiver AuthorizeUseCase) LoginInputDao(req request.UserLoginRequest) (response.LoginResponseData, error) {
	user := receiver.UserRepository.FindUserByUsername(req.Username)
	if user == nil {
		log.Info("No user has username matches", req.Username)
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

func (receiver AuthorizeUseCase) UserLoginUsecase(req request.UserLoginRequest) (response.LoginResponseData, error) {
	user, err := receiver.UserEntityRepository.GetByUsername(request.GetUserEntityByUsernameRequest{Username: req.Username})
	if err != nil {
		return response.LoginResponseData{}, errors.New("user not found")
	}

	if user == nil {
		log.Info("No user has username matches", req.Username)
		return response.LoginResponseData{}, errors.New("user not found")
	}

	err = receiver.SessionRepository.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return response.LoginResponseData{}, errors.New("invalid username or password")
	}

	token, err := receiver.SessionRepository.GenerateTokenV2(*user)
	if err != nil {
		return response.LoginResponseData{}, errors.New("cannot generate token")
	}

	//authMiddleware := jwtauth.JwtMiddleware()
	//token := authMiddleware.TokenGen(user.UserId)
	return *token, nil
}
