package usecase

import (
	"errors"
	"fmt"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	if tx.Error != nil {
		return fmt.Errorf("Failt create transaction: %s", tx.Error.Error())
	}

	// Đảm bảo rollback nếu có lỗi
	rolledBack := false
	defer func() {
		if !rolledBack {
			tx.Rollback()
		}
	}()

	// 1. Xoá dữ liệu cũ
	for _, item := range req {
		if err := receiver.ComponentRepository.DeleteBySectionID(item.SectionID, tx); err != nil {
			logrus.Error("Rollback error components section_id:", item.SectionID)
			tx.Rollback()
			rolledBack = true
			return fmt.Errorf("Delete components by section_id fail: %w", err)
		}
	}
	if err := receiver.ChildMenuRepository.DeleteAllTx(tx); err != nil {
		logrus.Error("Rollback error by delete child_menu:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("Delete child_menu fail: %w", err)
	}
	if err := receiver.StudentMenuRepository.DeleteAllTx(tx); err != nil {
		logrus.Error("Rollback error by delete student_menu:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("Delete student_menu fail: %w", err)
	}

	// 2. Lấy danh sách child_id và student_id
	childIDs, err := receiver.ChildRepository.GetAllIDs()
	if err != nil {
		logrus.Error("Rollback error by get child_ids:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("Get list child_id fail: %w", err)
	}

	studentIDs, err := receiver.StudentApplicationRepository.GetAllStudentIDs()
	if err != nil {
		logrus.Error("Rollback error by get student_ids:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("Get list student_id fail: %w", err)
	}

	// 3. Tạo component và gán menu theo Role
	for _, item := range req {
		parsedUUID, err := uuid.Parse(item.SectionID)
		if err != nil || parsedUUID == uuid.Nil {
			continue
		}

		roleOrg, err := receiver.RoleOrgSignUpRepository.GetByID(item.SectionID)
		if err != nil {
			logrus.Error("Rollback error get role by section_id :", err)
			tx.Rollback()
			rolledBack = true
			return fmt.Errorf("Get role by section_id fail: %w", err)
		}
		if roleOrg == nil {
			continue
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
				logrus.Error("Rollback by error create component:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Create component fail: %w", err)
			}

			visible, err := helper.GetVisibleToValueComponent(compReq.Value)
			if err != nil {
				logrus.Error("Rollback by error get visible:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Get Visible fail: %w", err)
			}

			switch roleOrg.RoleName {
			case string(value.RoleChild):
				if err := receiver.createChildMenus(tx, component.ID, visible, idx, childIDs); err != nil {
					logrus.Error("Rollback by error create child menu:", err)
					tx.Rollback()
					rolledBack = true
					return err
				}
			case string(value.RoleStudent):
				if err := receiver.createStudentMenus(tx, component.ID, visible, idx, studentIDs); err != nil {
					logrus.Error("Rollback by error create student menu:", err)
					tx.Rollback()
					rolledBack = true
					return err
				}
			}
		}
	}

	// 4. Commit transaction
	if err := tx.Commit().Error; err != nil {
		logrus.Error("Error commit transaction:", err)
		rolledBack = true
		return fmt.Errorf("commit transaction fail: %w", err)
	}

	rolledBack = true // Commit thành công, không rollback
	return nil
}

func (receiver *UploadSectionMenuUseCase) createChildMenus(tx *gorm.DB, componentID uuid.UUID, visible bool, order int, childIDs []uuid.UUID) error {
	for _, childID := range childIDs {
		existing, err := receiver.ChildMenuRepository.GetByChildIDAndComponentID(tx, childID, componentID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("get child menu fail: %w", err)
		}

		if existing != nil {
			// Đã tồn tại: cập nhật
			existing.Order = order
			existing.Visible = visible
			existing.IsShow = true
			if err := receiver.ChildMenuRepository.UpdateWithTx(tx, existing); err != nil {
				return fmt.Errorf("update child menu fail: %w", err)
			}
		} else {
			// Không tồn tại: tạo mới
			menu := &entity.ChildMenu{
				ID:          uuid.New(),
				ChildID:     childID,
				ComponentID: componentID,
				Order:       order,
				IsShow:      true,
				Visible:     visible,
			}
			if err := receiver.ChildMenuRepository.CreateWithTx(tx, menu); err != nil {
				return fmt.Errorf("create child menu fail: %w", err)
			}
		}
	}
	return nil
}

func (receiver *UploadSectionMenuUseCase) createStudentMenus(tx *gorm.DB, componentID uuid.UUID, visible bool, order int, studentIDs []uuid.UUID) error {
	for _, studentID := range studentIDs {
		existing, err := receiver.StudentMenuRepository.GetByStudentIDAndComponentID(tx, studentID, componentID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("lấy student menu thất bại: %w", err)
		}

		if existing != nil {
			// Đã tồn tại → update
			existing.Order = order
			existing.Visible = visible
			existing.IsShow = true
			if err := receiver.StudentMenuRepository.UpdateWithTx(tx, existing); err != nil {
				return fmt.Errorf("cập nhật student menu thất bại: %w", err)
			}
		} else {
			// Không tồn tại → create
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
	}
	return nil
}

func (receiver *UploadSectionMenuUseCase) UploadSectionMenuV2(req request.UploadSectionMenuRequest) error {
	tx := receiver.MenuRepository.DBConn.Begin()
	if tx.Error != nil {
		return fmt.Errorf("fail to create transaction: %s", tx.Error.Error())
	}

	rolledBack := false
	defer func() {
		if !rolledBack {
			tx.Rollback()
		}
	}()

	// 1. Lấy danh sách child_id và student_id
	childIDs, err := receiver.ChildRepository.GetAllIDs()
	if err != nil {
		logrus.Error("Rollback by error getting child_ids:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("Get list child_id failed: %w", err)
	}

	studentIDs, err := receiver.StudentApplicationRepository.GetAllStudentIDs()
	if err != nil {
		logrus.Error("Rollback by error getting student_ids:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("Get list student_id failed: %w", err)
	}

	// 2. Upsert component và tạo menu theo role
	for _, item := range req {

		// dau tien xoa component, child menu, student menu neu co mang delete_component_ids
		if len(item.DeleteComponentIDs) > 0 {
			for _, compID := range item.DeleteComponentIDs {
				if err := receiver.DeleteSectionMenu(compID); err != nil {
					logrus.Error("Rollback by error deleting section menu:", err)
					tx.Rollback()
					rolledBack = true
					return fmt.Errorf("Delete section menu failed: %w", err)
				}
			}
		}
		parsedUUID, err := uuid.Parse(item.SectionID)
		if err != nil || parsedUUID == uuid.Nil {
			continue
		}

		roleOrg, err := receiver.RoleOrgSignUpRepository.GetByID(item.SectionID)
		if err != nil {
			logrus.Error("Rollback by error getting role by section_id:", err)
			tx.Rollback()
			rolledBack = true
			return fmt.Errorf("Get role by section_id failed: %w", err)
		}
		if roleOrg == nil {
			continue
		}

		for idx, compReq := range item.Components {
			var componentID uuid.UUID

			component := &components.Component{
				Name:      compReq.Name,
				Type:      components.ComponentType(compReq.Type),
				Key:       compReq.Key,
				Value:     datatypes.JSON([]byte(compReq.Value)),
				SectionID: item.SectionID,
			}

			if compReq.ID != nil && *compReq.ID != uuid.Nil {
				// Nếu có ID truyền lên
				componentID = *compReq.ID
				component.ID = componentID

				existingComponent, err := receiver.ComponentRepository.GetByID(componentID.String())
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					logrus.Error("rollback by error query component:", err)
					tx.Rollback()
					rolledBack = true
					return fmt.Errorf("query component fail: %w", err)
				}

				if existingComponent != nil {
					// Update
					if err := receiver.ComponentRepository.UpdateWithTx(tx, component); err != nil {
						logrus.Error("rollback by error update component:", err)
						tx.Rollback()
						rolledBack = true
						return fmt.Errorf("update component fail: %w", err)
					}
				} else {
					// ID có nhưng không tồn tại
					component.ID = uuid.New()
					componentID = component.ID
					if err := receiver.ComponentRepository.CreateWithTx(tx, component); err != nil {
						logrus.Error("rollback by error create component (from non-existent id):", err)
						tx.Rollback()
						rolledBack = true
						return fmt.Errorf("create component fail: %w", err)
					}
				}
			} else {
				// Tạo mới
				component.ID = uuid.New()
				componentID = component.ID
				if err := receiver.ComponentRepository.CreateWithTx(tx, component); err != nil {
					logrus.Error("rollback by error create component:", err)
					tx.Rollback()
					rolledBack = true
					return fmt.Errorf("create component fail: %w", err)
				}
			}

			visible, err := helper.GetVisibleToValueComponent(compReq.Value)
			if err != nil {
				logrus.Error("rollback by error get visible:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("get visible fail: %w", err)
			}

			switch roleOrg.RoleName {
			case string(value.RoleChild):
				if err := receiver.createChildMenus(tx, componentID, visible, idx, childIDs); err != nil {
					logrus.Error("rollback by error create child menu:", err)
					tx.Rollback()
					rolledBack = true
					return err
				}
			case string(value.RoleStudent):
				if err := receiver.createStudentMenus(tx, componentID, visible, idx, studentIDs); err != nil {
					logrus.Error("rollback by error create student menu:", err)
					tx.Rollback()
					rolledBack = true
					return err
				}
			}
		}

	}

	if err := tx.Commit().Error; err != nil {
		logrus.Error("Error committing transaction:", err)
		rolledBack = true
		return fmt.Errorf("commit transaction failed: %s", err.Error())
	}

	rolledBack = true
	return nil
}

func (receiver *UploadSectionMenuUseCase) DeleteSectionMenu(componentID string) error {
	// Xóa component
	if err := receiver.ComponentRepository.DeleteComponent(componentID, nil); err != nil {
		return fmt.Errorf("UploadSectionMenuUseCase.DeleteSectionMenu: delete component failed: %w", err)
	}

	// Xóa child menu
	if err := receiver.ChildMenuRepository.DeleteByComponentID(componentID); err != nil {
		return fmt.Errorf("UploadSectionMenuUseCase.DeleteSectionMenu: delete child menu failed: %w", err)
	}

	// Xóa student menu
	if err := receiver.StudentMenuRepository.DeleteByComponentID(componentID); err != nil {
		return fmt.Errorf("UploadSectionMenuUseCase.DeleteSectionMenu: delete student menu failed: %w", err)
	}

	return nil
}
