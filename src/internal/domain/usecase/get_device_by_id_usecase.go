package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetDeviceByIdUseCase struct {
	*repository.DeviceRepository
}

func (receiver *GetDeviceByIdUseCase) Get(deviceId string) (*entity.SDevice, error) {
	return receiver.DeviceRepository.FindDeviceById(deviceId)
}
