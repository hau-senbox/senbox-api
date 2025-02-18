package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type GetDeviceComponentValuesUseCase struct {
	*repository.DeviceComponentValuesRepository
}

func (receiver *GetDeviceComponentValuesUseCase) GetDeviceComponentValuesByCompany(req request.GetDeviceComponentValuesByCompanyRequest) (*entity.SDeviceComponentValues, error) {
	return receiver.GetByCompany(req)
}

func (receiver *GetDeviceComponentValuesUseCase) GetDeviceComponentValuesByDevice(req request.GetDeviceComponentValuesByDeviceRequest) (*entity.SDeviceComponentValues, error) {
	return receiver.GetByDevice(req)
}
