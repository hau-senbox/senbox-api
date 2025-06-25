package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type TakeNoteUseCase struct {
	*repository.DeviceRepository
}

func (receiver *TakeNoteUseCase) TakeNote(params request.TakeNoteRequest, deviceID string) error {
	device, err := receiver.GetDeviceByID(deviceID)
	if err != nil {
		return err
	}

	device.Note = params.Note
	_, err = receiver.UpdateDevice(device)
	if err != nil {
		return err
	}
	return nil
}
