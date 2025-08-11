package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetOrganizationUseCase struct {
	*repository.OrganizationRepository
	*repository.DeviceRepository
}

func (receiver *GetOrganizationUseCase) GetOrganizationByID(id string) (*entity.SOrganization, error) {
	return receiver.OrganizationRepository.GetByID(id)
}

func (receiver *GetOrganizationUseCase) GetByName(name string) (*entity.SOrganization, error) {
	return receiver.OrganizationRepository.GetByName(name)
}

func (receiver *GetOrganizationUseCase) GetAllOrganization(user *entity.SUserEntity) ([]*entity.SOrganization, error) {
	return receiver.OrganizationRepository.GetAll(user)
}

func (receiver *GetOrganizationUseCase) GetAllUserByOrganization(organizationID string) ([]*entity.SUserOrg, error) {
	return receiver.OrganizationRepository.GetAllUserByOrganization(organizationID)
}

func (receiver *GetOrganizationUseCase) CheckDeviceInOrg4App(deviceID string, organizationID string) (bool, error) {
	deviceOrgIds, _ := receiver.DeviceRepository.GetOrgIDsByDeviceID(deviceID)

	found := false
	for _, orgID := range deviceOrgIds {
		if orgID.String() == organizationID {
			found = true
			break
		}
	}

	return found, nil
}
