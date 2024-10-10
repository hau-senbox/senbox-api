package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type TakeNoteUseCase struct {
	*repository.DeviceRepository
}

func (receiver *TakeNoteUseCase) TakeNote(params request.TakeNoteRequest, device *entity.SDevice) error {
	device.Note = params.Note
	_, err := receiver.DeviceRepository.UpdateDevice(device)
	if err != nil {
		return err
	}
	return nil
}
