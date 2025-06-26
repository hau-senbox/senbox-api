package usecase

import (
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity/menu"
	"sen-global-api/internal/domain/request"
)

type UploadSuperAdminMenuUseCase struct {
	*repository.MenuRepository
	*repository.ComponentRepository
}

func (receiver *UploadSuperAdminMenuUseCase) Upload(req request.UploadSuperAdminMenuRequest) error {
	tx := receiver.MenuRepository.DBConn.Begin()
	if err := receiver.MenuRepository.DeleteSuperAdminMenu(tx); err != nil {
		return err
	}

	if len(req.Top) > 0 {
		err := receiver.ComponentRepository.CreateComponents(&req.Top, tx)
		if err != nil {
			return err
		}

		if err := receiver.MenuRepository.CreateSuperAdminMenu(request.CreateSuperAdminMenuRequest{
			Direction:  menu.Top,
			Components: req.Top,
		}, tx); err != nil {
			return fmt.Errorf("failed to create top menu: %w", err)
		}
	}

	if len(req.Bottom) > 0 {
		err := receiver.ComponentRepository.CreateComponents(&req.Bottom, tx)
		if err != nil {
			return err
		}

		if err := receiver.MenuRepository.CreateSuperAdminMenu(request.CreateSuperAdminMenuRequest{
			Direction:  menu.Bottom,
			Components: req.Bottom,
		}, tx); err != nil {
			return fmt.Errorf("failed to create bottom menu: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
