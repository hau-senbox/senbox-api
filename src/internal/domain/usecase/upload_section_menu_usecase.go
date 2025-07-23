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
		return fmt.Errorf("khởi tạo transaction thất bại: %w", tx.Error)
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
			logrus.Error("Rollback do lỗi xóa components section_id:", item.SectionID)
			tx.Rollback()
			rolledBack = true
			return fmt.Errorf("xóa components theo section_id thất bại: %w", err)
		}
	}
	if err := receiver.ChildMenuRepository.DeleteAllTx(tx); err != nil {
		logrus.Error("Rollback do lỗi xóa child_menu:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("xóa child_menu thất bại: %w", err)
	}
	if err := receiver.StudentMenuRepository.DeleteAllTx(tx); err != nil {
		logrus.Error("Rollback do lỗi xóa student_menu:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("xóa student_menu thất bại: %w", err)
	}

	// 2. Lấy danh sách child_id và student_id
	childIDs, err := receiver.ChildRepository.GetAllIDs()
	if err != nil {
		logrus.Error("Rollback do lỗi lấy child_ids:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("lấy danh sách child_id thất bại: %w", err)
	}

	studentIDs, err := receiver.StudentApplicationRepository.GetAllStudentIDs()
	if err != nil {
		logrus.Error("Rollback do lỗi lấy student_ids:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("lấy danh sách student_id thất bại: %w", err)
	}

	// 3. Tạo component và gán menu theo Role
	for _, item := range req {
		parsedUUID, err := uuid.Parse(item.SectionID)
		if err != nil || parsedUUID == uuid.Nil {
			continue
		}

		roleOrg, err := receiver.RoleOrgSignUpRepository.GetByID(item.SectionID)
		if err != nil {
			logrus.Error("Rollback do lỗi lấy role theo section_id :", err)
			tx.Rollback()
			rolledBack = true
			return fmt.Errorf("lấy role theo section_id thất bại: %w", err)
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
				logrus.Error("Rollback do lỗi tạo component:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("tạo component thất bại: %w", err)
			}

			visible, err := helper.GetVisibleToValueComponent(compReq.Value)
			if err != nil {
				logrus.Error("Rollback do lỗi phân tích visible:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("phân tích Visible thất bại: %w", err)
			}

			switch roleOrg.RoleName {
			case string(value.RoleChild):
				if err := receiver.createChildMenus(tx, component.ID, visible, idx, childIDs); err != nil {
					logrus.Error("Rollback do lỗi tạo child menu:", err)
					tx.Rollback()
					rolledBack = true
					return err
				}
			case string(value.RoleStudent):
				if err := receiver.createStudentMenus(tx, component.ID, visible, idx, studentIDs); err != nil {
					logrus.Error("Rollback do lỗi tạo student menu:", err)
					tx.Rollback()
					rolledBack = true
					return err
				}
			}
		}
	}

	// 4. Commit transaction
	if err := tx.Commit().Error; err != nil {
		logrus.Error("Lỗi commit transaction:", err)
		rolledBack = true
		return fmt.Errorf("commit transaction thất bại: %w", err)
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
