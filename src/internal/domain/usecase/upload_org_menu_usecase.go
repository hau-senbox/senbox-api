package usecase

import (
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity/menu"
	"sen-global-api/internal/domain/request"
)

type UploadOrgMenuUseCase struct {
	*repository.MenuRepository
	*repository.ComponentRepository
}

func (receiver *UploadOrgMenuUseCase) Upload(req request.UploadOrgMenuRequest) error {
	tx := receiver.MenuRepository.DBConn.Begin()
	if err := receiver.MenuRepository.DeleteOrgMenu(req.OrganizationID, tx); err != nil {
		return err
	}

	if len(req.Top) > 0 {
		err := receiver.ComponentRepository.CreateComponents(&req.Top, tx)
		if err != nil {
			return err
		}

		if err := receiver.MenuRepository.CreateOrgMenu(request.CreateOrgMenuRequest{
			OrganizationID: req.OrganizationID,
			Direction:      menu.Top,
			Components:     req.Top,
		}, tx); err != nil {
			return fmt.Errorf("failed to create top menu: %w", err)
		}
	}

	if len(req.Bottom) > 0 {
		err := receiver.ComponentRepository.CreateComponents(&req.Bottom, tx)
		if err != nil {
			return err
		}

		if err := receiver.MenuRepository.CreateOrgMenu(request.CreateOrgMenuRequest{
			OrganizationID: req.OrganizationID,
			Direction:      menu.Bottom,
			Components:     req.Bottom,
		}, tx); err != nil {
			return fmt.Errorf("failed to create bottom menu: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
