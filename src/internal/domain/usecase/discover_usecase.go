package usecase

import (
	"sen-global-api/internal/data/repository"
)

type DiscoverUseCase struct {
	*repository.DeviceRepository
}

type DiscoveredDeviceData struct {
	DeviceID string `json:"device_id"`
}

func (receiver *DiscoverUseCase) Execute(deviceID string) (DiscoveredDeviceData, error) {
	device, err := receiver.GetDeviceByID(deviceID)
	if err != nil {
		return DiscoveredDeviceData{}, err
	}

	return DiscoveredDeviceData{
		DeviceID: device.ID,
	}, nil
}
