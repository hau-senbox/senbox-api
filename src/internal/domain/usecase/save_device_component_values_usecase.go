package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type SaveDeviceComponentValuesUseCase struct {
	*repository.DeviceComponentValuesRepository
}

func (receiver *SaveDeviceComponentValuesUseCase) SaveDeviceComponentValuesByOrganization(req request.SaveDeviceComponentValuesByOrganizationRequest) error {
	return receiver.SaveByOrganization(req)
}

func (receiver *SaveDeviceComponentValuesUseCase) SaveDeviceComponentValuesByDevice(req request.SaveDeviceComponentValuesByDeviceRequest) error {
	return receiver.SaveByDevice(req)
}
