package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type GetUserEntityUseCase struct {
	*repository.UserEntityRepository
	*repository.OrganizationRepository
}

func (receiver *GetUserEntityUseCase) GetUserByID(req request.GetUserEntityByIDRequest) (*entity.SUserEntity, error) {
	return receiver.UserEntityRepository.GetByID(req)
}

func (receiver *GetUserEntityUseCase) GetUserByUsername(req request.GetUserEntityByUsernameRequest) (*entity.SUserEntity, error) {
	return receiver.UserEntityRepository.GetByUsername(req)
}

func (receiver *GetUserEntityUseCase) GetAllUsers() ([]entity.SUserEntity, error) {
	return receiver.UserEntityRepository.GetAll()
}

func (receiver *GetUserEntityUseCase) GetAllByOrganization(organizationID string) ([]entity.SUserEntity, error) {
	return receiver.UserEntityRepository.GetAllByOrganizationID(organizationID)
}

func (receiver *GetUserEntityUseCase) GetUserOrgInfo(userID, organization string) (*entity.SUserOrg, error) {
	return receiver.OrganizationRepository.GetUserOrgInfo(userID, organization)
}

func (receiver *GetUserEntityUseCase) GetAllOrgManagerInfo(organization string) (*[]entity.SUserOrg, error) {
	return receiver.OrganizationRepository.GetAllOrgManagerInfo(organization)
}

func (receiver *GetUserEntityUseCase) GetAllUserAuthorize(userID string) ([]entity.SUserFunctionAuthorize, error) {
	return receiver.UserEntityRepository.GetAllUserAuthorize(userID)
}
