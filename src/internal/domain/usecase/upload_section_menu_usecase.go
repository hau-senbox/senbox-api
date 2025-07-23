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
		return fmt.Errorf("Failt create transaction: %w", tx.Error.Error())
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
