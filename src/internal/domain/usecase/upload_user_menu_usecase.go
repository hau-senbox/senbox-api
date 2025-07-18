package usecase

import (
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/request"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UploadUserMenuUseCase struct {
	*repository.MenuRepository
	*repository.ComponentRepository
	*repository.ChildMenuRepository
	*repository.ChildRepository
}

func (receiver *UploadUserMenuUseCase) Upload(req request.UploadUserMenuRequest) error {
	tx := receiver.MenuRepository.DBConn.Begin()
	if err := receiver.MenuRepository.DeleteUserMenu(req.UserID, tx); err != nil {
		return err
	}

	if len(req.Components) > 0 {
		err := receiver.ComponentRepository.CreateComponents(&req.Components, tx)
		if err != nil {
			return err
		}

		if err := receiver.MenuRepository.CreateUserMenu(request.CreateUserMenuRequest{
			UserID:     req.UserID,
			Components: req.Components,
		}, tx); err != nil {
			return fmt.Errorf("failed to create user menu: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (receiver *UploadUserMenuUseCase) UploadSectionMenu(req request.UploadSectionMenuRequest) error {
	tx := receiver.MenuRepository.DBConn.Begin()

	for _, item := range req {
		// Xoá tất cả components theo section_id
		if err := receiver.ComponentRepository.DeleteBySectionID(item.SectionID, tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete components by section: %w", err)
		}
		// Xoa trong child_menu
		if err := receiver.ChildMenuRepository.DeleteAll(); err != nil {
			return fmt.Errorf("failed to delete child menu: %w", err)
		}
	}

	// Lấy danh sách tất cả child_id
	childIDs, err := receiver.ChildRepository.GetAllIDs()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get child IDs: %w", err)
	}

	// Thêm mới components
	for _, item := range req {
		for index, compReq := range item.Components {
			component := &components.Component{
				ID:        uuid.New(),
				Name:      compReq.Name,
				Type:      components.ComponentType(compReq.Type),
				Key:       compReq.Key,
				Value:     datatypes.JSON([]byte(compReq.Value)),
				SectionID: item.SectionID,
			}

			if err := receiver.ComponentRepository.CreateWithTx(tx, component); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create component: %w", err)
			}

			// Gắn component với từng child_id
			for _, childID := range childIDs {
				childMenu := &entity.ChildMenu{
					ID:          uuid.New(),
					ChildID:     childID,
					ComponentID: component.ID,
					Order:       index,
					IsShow:      true,
				}
				if err := receiver.ChildMenuRepository.CreateWithTx(tx, childMenu); err != nil {
					tx.Rollback()
					return fmt.Errorf("failed to create child menu: %w", err)
				}
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
