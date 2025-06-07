package usecase

import (
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity/menu"
	"sen-global-api/internal/domain/request"
)

type UploadUserMenuUseCase struct {
	*repository.MenuRepository
	*repository.ComponentRepository
}

func (receiver *UploadUserMenuUseCase) Upload(req request.UploadUserMenuRequest) error {
	tx := receiver.MenuRepository.DBConn.Begin()
	if err := receiver.MenuRepository.DeleteUserMenu(req.UserID, tx); err != nil {
		return err
	}

	err := receiver.ComponentRepository.CreateComponents(&req.Components, tx)
	if err != nil {
		return err
	}

	if err := receiver.MenuRepository.CreateUserMenu(request.CreateUserMenuRequest{
		UserID:     req.UserID,
		Direction:  menu.Top,
		Components: req.Components,
	}, tx); err != nil {
		return fmt.Errorf("failed to create user menu: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
