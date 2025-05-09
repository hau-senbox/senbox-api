package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/value"
)

type RefreshAccessTokenUseCase struct {
	*repository.SessionRepository
	*repository.DeviceRepository
}

func (c *RefreshAccessTokenUseCase) Execute(refreshToken string) (string, string, error) {
	deviceID, err := c.GetDeviceIDFromRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	device, err := c.FindDeviceById(deviceID)
	if err != nil || device == nil {
		return "", "", err
	}

	if device.Status != value.DeviceModeSuspended {
		return "", "", err
	}

	accessToken, refreshToken, err := c.GenerateTokenByDevice(*device)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
