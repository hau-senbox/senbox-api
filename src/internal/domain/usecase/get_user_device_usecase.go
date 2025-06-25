package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetUserDeviceUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *GetUserDeviceUseCase) GetUserDeviceByID(deviceID string) (*[]entity.SUserDevices, error) {
	return receiver.UserEntityRepository.GetUserDeviceByID(deviceID)
}
