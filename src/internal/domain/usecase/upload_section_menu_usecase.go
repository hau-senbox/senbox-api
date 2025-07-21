package usecase

import (
	"fmt"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UploadSectionMenuUseCase struct {
	*repository.MenuRepository
	*repository.ComponentRepository
	*repository.ChildMenuRepository
	*repository.ChildRepository
	*repository.RoleOrgSignUpRepository
	*repository.StudentMenuRepository
	*repository.StudentApplicationRepository
}

func (receiver *UploadSectionMenuUseCase) UploadSectionMenu(req request.UploadSectionMenuRequest) error {
	tx := receiver.MenuRepository.DBConn.Begin()

	for _, item := range req {
		// Xoá tất cả components theo section_id
		if err := receiver.ComponentRepository.DeleteBySectionID(item.SectionID, tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete components by section: %w", err)
		}
		// Xoá toàn bộ child_menu (nếu muốn giới hạn theo section_id thì cần chỉnh lại repository)
		if err := receiver.ChildMenuRepository.DeleteAll(); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete child menu: %w", err)
		}
	}

	// Lấy danh sách tất cả child_id
	childIDs, err := receiver.ChildRepository.GetAllIDs()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get child IDs: %w", err)
	}

	// Lay danh sach student
	studentIDs, err := receiver.StudentApplicationRepository.GetAllStudentIDs()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get student IDs: %w", err)
	}

	// Tạo mới component và gắn vào child nếu cần
	for _, item := range req {
		// Lấy role theo SectionID
		roleOrg, err := receiver.RoleOrgSignUpRepository.GetByID(item.SectionID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get role org by ID: %w", err)
		}

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
				return fmt.Errorf("failed to get visible value: %w", err)
			}

			// Nếu là role "Child" thì gắn vào bảng ChildMenu
			if roleOrg != nil && roleOrg.RoleName == string(value.RoleChild) {
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

			// Nếu là role "Student" thì gắn vào bảng StudentMenu
			if roleOrg != nil && roleOrg.RoleName == string(value.RoleStudent) {
				for _, studentID := range studentIDs {
					studentMenu := &entity.StudentMenu{
						ID:          uuid.New(),
						StudentID:   studentID,
						ComponentID: component.ID,
						Order:       index,
						IsShow:      true,
						Visible:     visible,
					}
					if err := receiver.StudentMenuRepository.CreateWithTx(tx, studentMenu); err != nil {
						tx.Rollback()
						return fmt.Errorf("failed to create student menu: %w", err)
					}
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
