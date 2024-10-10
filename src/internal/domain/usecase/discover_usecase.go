package usecase

import (
	"sen-global-api/internal/data/repository"
)

type DiscoverUseCase struct {
	*repository.DeviceRepository
}

type DiscoveredDeviceData struct {
	DeviceId  string `json:"device_id"`
	UserInfo1 string `json:"user_info_1"`
	UserInfo2 string `json:"user_info_2"`
	UserInfo3 string `json:"user_info_3"`
}

func (receiver *DiscoverUseCase) Execute(deviceId string) (DiscoveredDeviceData, error) {
	device, err := receiver.DeviceRepository.GetDeviceById(deviceId)
	if err != nil {
		return DiscoveredDeviceData{}, err
	}

	return DiscoveredDeviceData{
		DeviceId:  device.DeviceId,
		UserInfo1: device.PrimaryUserInfo,
		UserInfo2: device.SecondaryUserInfo,
		UserInfo3: device.TertiaryUserInfo,
	}, nil
}
