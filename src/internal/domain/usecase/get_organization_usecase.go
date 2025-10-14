package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
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

func (receiver *GetOrganizationUseCase) GetAllOrganizations4Gateway() ([]response.OrganizationResponse, error) {
	orgs, err := receiver.OrganizationRepository.GetAll4Gateway()
	if err != nil {
		return nil, err
	}

	// map sang response
	responses := make([]response.OrganizationResponse, 0, len(orgs))
	for _, org := range orgs {
		responses = append(responses, response.OrganizationResponse{
			ID:               org.ID.String(),
			OrganizationName: org.OrganizationName,
			Avatar:           org.Avatar,
			AvatarURL:        org.AvatarURL,
			Address:          org.Address,
			Description:      org.Description,
		})
	}

	return responses, nil
}
