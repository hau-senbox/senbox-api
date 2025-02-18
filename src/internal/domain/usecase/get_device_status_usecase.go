package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
)

type GetDeviceStatusUseCase struct {
	DeviceRepository *repository.DeviceRepository
}

func statusInStringFrom(status value.DeviceMode) string {
	switch status {
	case value.DeviceModeSuspended:
		return "suspended"
	case value.DeviceModeT:
		return "mode_t"
	case value.DeviceModeS:
		return "mode_s"
	case value.DeviceModeP:
		return "mode_p"
	case value.DeviceModeDeactivated:
		return "deactivated"
	case value.DeviceModeL:
		return "mode_l"
	}
	return ""
}

func (receiver *GetDeviceStatusUseCase) Execute(device entity.SDevice) (response.GetDeviceStatusResponseData, error) {
	return response.GetDeviceStatusResponseData{
		Status:  statusInStringFrom(device.Status),
		Message: device.DeactivateMessage,
	}, nil
}
