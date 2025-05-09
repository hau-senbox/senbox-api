package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
)

type GetUserEntityUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *GetUserEntityUseCase) GetUserById(req request.GetUserEntityByIdRequest) (*entity.SUserEntity, error) {
	return receiver.GetByID(req)
}

func (receiver *GetUserEntityUseCase) GetChildrenOfGuardian(userId string) (*[]response.UserEntityResponseData, error) {
	return receiver.UserEntityRepository.GetChildrenOfGuardian(userId)
}

func (receiver *GetUserEntityUseCase) GetUserByUsername(req request.GetUserEntityByUsernameRequest) (*entity.SUserEntity, error) {
	return receiver.GetByUsername(req)
}

func (receiver *GetUserEntityUseCase) GetAllUsers() ([]entity.SUserEntity, error) {
	return receiver.GetAll()
}

func (receiver *GetUserEntityUseCase) GetUserOrgInfo(userId, organization string) (*entity.SUserOrg, error) {
	return receiver.UserEntityRepository.GetUserOrgInfo(userId, organization)
}

func (receiver *GetUserEntityUseCase) GetAllOrgManagerInfo(organization string) (*[]entity.SUserOrg, error) {
	return receiver.UserEntityRepository.GetAllOrgManagerInfo(organization)
}

func (receiver *GetUserEntityUseCase) GetAllUserAuthorize(userId string) ([]entity.SUserFunctionAuthorize, error) {
	return receiver.UserEntityRepository.GetAllUserAuthorize(userId)
}
