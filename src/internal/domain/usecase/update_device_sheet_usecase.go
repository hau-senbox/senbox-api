package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UpdateDeviceSheetUseCase struct {
	DeviceRepository *repository.DeviceRepository
}

func (receiver UpdateDeviceSheetUseCase) DeactivateDevice(deviceID string, req request.DeactivateDeviceRequest) error {
	return receiver.DeviceRepository.DeactivateDevice(deviceID, req.Message)
}

func (receiver UpdateDeviceSheetUseCase) ActivateDevice(deviceID string, req request.ReactivateDeviceRequest) error {
	return receiver.DeviceRepository.ActivateDevice(deviceID, req.Message)
}
