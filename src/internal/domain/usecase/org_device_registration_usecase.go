package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
)

type OrgDeviceRegistrationUseCase struct {
	*repository.OrganizationRepository
	*repository.DeviceRepository
}

func (receiver *OrgDeviceRegistrationUseCase) RegisterOrgDevice(req request.RegisteringDeviceForOrg) error {
	org, err := receiver.OrganizationRepository.GetByID(req.OrgID)
	if err != nil {
		return err
	}

	if err := receiver.CheckOrgDeviceExist(req); err != nil {
		return err
	}

	_, err = receiver.RegisteringDeviceForOrg(org, request.RegisterDeviceRequest{
		DeviceUUID: req.DeviceID,
		InputMode:  string(value.InfoInputTypeBarcode),
	})
	if err != nil {
		return err
	}

	return nil
}
