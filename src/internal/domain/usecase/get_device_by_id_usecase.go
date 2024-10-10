package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type GetDeviceByIdUseCase struct {
	*repository.DeviceRepository
}

func (receiver *GetDeviceByIdUseCase) Get(req request.ReconnectDeviceRequest) (*entity.SDevice, error) {
	return receiver.DeviceRepository.FindDeviceById(req.DeviceId)
}
