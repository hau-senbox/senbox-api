package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type UpdateUserEntityUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *UpdateUserEntityUseCase) UpdateUserEntity(req request.UpdateUserEntityRequest) error {
	return receiver.UpdateUser(req)
}

func (receiver *UpdateUserEntityUseCase) BlockUser(userID string) error {
	return receiver.UserEntityRepository.BlockUser(userID)
}

func (uc *UpdateUserEntityUseCase) UpdateCustomIDByUserID(req request.AddCustomID2UserRequest) error {
	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error("UpdateCustomIDByUserID: invalid user ID - " + err.Error())
		return errors.New("invalid user ID")
	}

	return uc.UserEntityRepository.UpdateCustomIDByUserID(userUUID, req.CustomID)
}

func (uc *UpdateUserEntityUseCase) SetReLogin() error {
	return uc.UserEntityRepository.SetReLogin()
}

func (uc *UpdateUserEntityUseCase) UpdateReLogin(userID string, relogin bool) error {
	return uc.UserEntityRepository.UpdateReLogin(userID, &relogin)
}
