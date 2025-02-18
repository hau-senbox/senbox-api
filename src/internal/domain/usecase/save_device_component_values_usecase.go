package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type SaveDeviceComponentValuesUseCase struct {
	*repository.DeviceComponentValuesRepository
}

func (receiver *SaveDeviceComponentValuesUseCase) SaveDeviceComponentValuesByCompany(req request.SaveDeviceComponentValuesByCompanyRequest) error {
	return receiver.SaveByCompany(req)
}

func (receiver *SaveDeviceComponentValuesUseCase) SaveDeviceComponentValuesByDevice(req request.SaveDeviceComponentValuesByDeviceRequest) error {
	return receiver.SaveByDevice(req)
}
