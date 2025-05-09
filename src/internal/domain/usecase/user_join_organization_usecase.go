package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserJoinOrganizationUseCase struct {
	*repository.OrganizationRepository
	repository.SessionRepository
}

func (receiver *UserJoinOrganizationUseCase) UserJoinOrganization(req request.UserJoinOrganizationRequest) error {
	organization, err := receiver.GetByID(req.OrganizationId)

	if err != nil {
		log.Error("OrganizationRepository.UserJoinOrganization: " + err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("organization doesn't exist")
		}
		return errors.New("failed to get organization")
	}

	err = receiver.VerifyPassword(req.Password, organization.Password)
	if err != nil {
		return errors.New("invalid organization or password")
	}

	return receiver.OrganizationRepository.UserJoinOrganization(req)
}
