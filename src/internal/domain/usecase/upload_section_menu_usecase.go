package usecase

import (
	"fmt"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/request"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UploadSectionMenuUseCase struct {
	*repository.MenuRepository
	*repository.ComponentRepository
	*repository.ChildMenuRepository
	*repository.ChildRepository
}

func (receiver *UploadSectionMenuUseCase) UploadSectionMenu(req request.UploadSectionMenuRequest) error {
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

			visible, err := helper.GetVisibleToValueComponent(compReq.Value)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Gắn component với từng child_id
			for _, childID := range childIDs {
				childMenu := &entity.ChildMenu{
					ID:          uuid.New(),
					ChildID:     childID,
					ComponentID: component.ID,
					Order:       index,
					IsShow:      true,
					Visible:     visible,
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
