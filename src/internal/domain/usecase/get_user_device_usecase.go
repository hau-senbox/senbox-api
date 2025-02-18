package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetUserDeviceUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *GetUserDeviceUseCase) GetUserDeviceById(deviceId string) (*[]entity.SUserDevices, error) {
	return receiver.UserEntityRepository.GetUserDeviceById(deviceId)
}
