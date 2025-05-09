package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type GetDeviceComponentValuesUseCase struct {
	*repository.DeviceComponentValuesRepository
}

func (receiver *GetDeviceComponentValuesUseCase) GetDeviceComponentValuesByOrganization(req request.GetDeviceComponentValuesByOrganizationRequest) (*entity.SDeviceComponentValues, error) {
	return receiver.GetByOrganization(req)
}

func (receiver *GetDeviceComponentValuesUseCase) GetDeviceComponentValuesByDevice(req request.GetDeviceComponentValuesByDeviceRequest) (*entity.SDeviceComponentValues, error) {
	return receiver.GetByDevice(req)
}
