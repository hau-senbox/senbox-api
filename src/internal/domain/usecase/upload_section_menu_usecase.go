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
	"gorm.io/gorm"
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

	// 1. Xoá dữ liệu cũ: Component, ChildMenu, StudentMenu
	for _, item := range req {
		if err := receiver.ComponentRepository.DeleteBySectionID(item.SectionID, tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("xóa components theo section_id thất bại: %w", err)
		}
	}
	if err := receiver.ChildMenuRepository.DeleteAll(); err != nil {
		tx.Rollback()
		return fmt.Errorf("xóa child_menu thất bại: %w", err)
	}
	if err := receiver.StudentMenuRepository.DeleteAll(); err != nil {
		tx.Rollback()
		return fmt.Errorf("xóa student_menu thất bại: %w", err)
	}

	// 2. Lấy danh sách child_id và student_id
	childIDs, err := receiver.ChildRepository.GetAllIDs()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("lấy danh sách child_id thất bại: %w", err)
	}

	studentIDs, err := receiver.StudentApplicationRepository.GetAllStudentIDs()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("lấy danh sách student_id thất bại: %w", err)
	}

	// 3. Tạo component và gắn vào menu tương ứng
	for _, item := range req {
		roleOrg, err := receiver.RoleOrgSignUpRepository.GetByID(item.SectionID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("lấy role theo section_id thất bại: %w", err)
		}
		if roleOrg == nil {
			continue // không có role -> bỏ qua
		}

		for idx, compReq := range item.Components {
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
				return fmt.Errorf("tạo component thất bại: %w", err)
			}

			visible, err := helper.GetVisibleToValueComponent(compReq.Value)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("phân tích Visible thất bại: %w", err)
			}

			switch roleOrg.RoleName {
			case string(value.RoleChild):
				if err := receiver.createChildMenus(tx, component.ID, visible, idx, childIDs); err != nil {
					tx.Rollback()
					return err
				}

			case string(value.RoleStudent):
				if err := receiver.createStudentMenus(tx, component.ID, visible, idx, studentIDs); err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	// 4. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction thất bại: %w", err)
	}

	return nil
}

func (receiver *UploadSectionMenuUseCase) createChildMenus(tx *gorm.DB, componentID uuid.UUID, visible bool, order int, childIDs []uuid.UUID) error {
	for _, childID := range childIDs {
		menu := &entity.ChildMenu{
			ID:          uuid.New(),
			ChildID:     childID,
			ComponentID: componentID,
			Order:       order,
			IsShow:      true,
			Visible:     visible,
		}
		if err := receiver.ChildMenuRepository.CreateWithTx(tx, menu); err != nil {
			return fmt.Errorf("tạo child menu thất bại: %w", err)
		}
	}
	return nil
}

func (receiver *UploadSectionMenuUseCase) createStudentMenus(tx *gorm.DB, componentID uuid.UUID, visible bool, order int, studentIDs []uuid.UUID) error {
	for _, studentID := range studentIDs {
		menu := &entity.StudentMenu{
			ID:          uuid.New(),
			StudentID:   studentID,
			ComponentID: componentID,
			Order:       order,
			IsShow:      true,
			Visible:     visible,
		}
		if err := receiver.StudentMenuRepository.CreateWithTx(tx, menu); err != nil {
			return fmt.Errorf("tạo student menu thất bại: %w", err)
		}
	}
	return nil
}
