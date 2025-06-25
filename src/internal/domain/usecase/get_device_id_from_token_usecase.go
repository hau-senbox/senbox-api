package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetDeviceIDFromTokenUseCase struct {
	*repository.SessionRepository
	*repository.DeviceRepository
}

func (receiver *GetDeviceIDFromTokenUseCase) GetDeviceFromToken(tokenString string) (*entity.SDevice, error) {
	deviceID, err := receiver.ExtractDeviceIDFromToken(tokenString)
	if err != nil {
		return nil, err
	}

	return receiver.FindDeviceByID(*deviceID)
}
