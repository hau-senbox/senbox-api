package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetDevicesByUserIdUseCase struct {
	*repository.DeviceRepository
}

func (receiver *GetDevicesByUserIdUseCase) GetDevicesByUserId(userId string) (*[]entity.SDevice, error) {
	return receiver.DeviceRepository.GetDevicesByUserId(userId)
}
