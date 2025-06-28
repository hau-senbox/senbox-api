package usecase

import (
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UploadDeviceMenuUseCase struct {
	*repository.MenuRepository
	*repository.ComponentRepository
}

func (receiver *UploadDeviceMenuUseCase) Upload(req request.UploadDeviceMenuRequest) error {
	tx := receiver.MenuRepository.DBConn.Begin()
	if err := receiver.MenuRepository.DeleteDeviceMenu(req.DeviceID, req.OrganizationID, tx); err != nil {
		return err
	}

	if len(req.Components) > 0 {
		err := receiver.ComponentRepository.CreateComponents(&req.Components, tx)
		if err != nil {
			return err
		}

		if err := receiver.MenuRepository.CreateDeviceMenu(request.CreateDeviceMenuRequest{
			DeviceID:       req.DeviceID,
			OrganizationID: req.OrganizationID,
			Components:     req.Components,
		}, tx); err != nil {
			return fmt.Errorf("failed to create device menu: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
