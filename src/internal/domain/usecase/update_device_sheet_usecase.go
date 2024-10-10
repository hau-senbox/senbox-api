package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UpdateDeviceSheetUseCase struct {
	DeviceRepository *repository.DeviceRepository
}

func (receiver UpdateDeviceSheetUseCase) DeactivateDevice(deviceId string, req request.DeactivateDeviceRequest) error {
	return receiver.DeviceRepository.DeactivateDevice(deviceId, req.Message)
}

func (receiver UpdateDeviceSheetUseCase) ActivateDevice(deviceId string, req request.ReactivateDeviceRequest) error {
	return receiver.DeviceRepository.ActivateDevice(deviceId, req.Message)
}
