package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetDeviceIdFromTokenUseCase struct {
	*repository.SessionRepository
	*repository.DeviceRepository
}

func (receiver *GetDeviceIdFromTokenUseCase) GetDeviceFromToken(tokenString string) (*entity.SDevice, error) {
	deviceId, err := receiver.ExtractDeviceIdFromToken(tokenString)
	if err != nil {
		return nil, err
	}

	return receiver.FindDeviceById(*deviceId)
}
