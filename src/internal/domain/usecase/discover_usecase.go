package usecase

import (
	"sen-global-api/internal/data/repository"
)

type DiscoverUseCase struct {
	*repository.DeviceRepository
}

type DiscoveredDeviceData struct {
	DeviceId string `json:"device_id"`
}

func (receiver *DiscoverUseCase) Execute(deviceId string) (DiscoveredDeviceData, error) {
	device, err := receiver.GetDeviceById(deviceId)
	if err != nil {
		return DiscoveredDeviceData{}, err
	}

	return DiscoveredDeviceData{
		DeviceId: device.ID,
	}, nil
}
