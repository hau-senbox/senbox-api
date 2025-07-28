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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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
	*GetUserEntityUseCase
	*repository.OrganizationMenuTemplateRepository
	*repository.TeacherApplicationRepository
	*repository.TeacherMenuRepository
	*repository.StaffMenuRepository
	*repository.StaffApplicationRepository
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

func (receiver *UploadSectionMenuUseCase) createStudentMenus(ctx *gin.Context, tx *gorm.DB, componentID uuid.UUID, visible bool, order int, studentIDs []uuid.UUID, sectionID uuid.UUID) error {
	// dau tien kiem tra user dang la quan ly cua organization nao
	user, err := receiver.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return err
	}

	// Nếu không phải SuperAdmin → lấy danh sách org mà user quản lý
	if len(user.Organizations) == 0 {
		return errors.New("user does not belong to any organization")
	}

	orgIDsManaged, err := user.GetManagedOrganizationIDs(receiver.UserEntityRepository.GetDB())
	if err != nil {
		return err
	}
	if len(orgIDsManaged) == 0 {
		return errors.New("user does not manage any organization")
	}

	// Tạo hoặc update OrganizationMenuTemplate cho mỗi tổ chức quản lý
	for _, orgID := range orgIDsManaged {
		existingTemplate, err := receiver.OrganizationMenuTemplateRepository.GetByOrgIDComponentIDSectionID(
			tx,
			orgID,
			componentID,
			sectionID,
		)
		if err != nil {
			log.Errorf("Error check OrganizationMenuTemplate: %v", err)
			return fmt.Errorf("check OrganizationMenuTemplate fail: %w", err)
		}

		if existingTemplate == nil {
			newTemplate := &entity.OrganizationMenuTemplate{
				ID:             uuid.New().String(),
				OrganizationID: orgID,
				ComponentID:    componentID.String(),
				SectionID:      sectionID.String(),
			}
			if err := receiver.OrganizationMenuTemplateRepository.CreateWithTx(tx, newTemplate); err != nil {
				log.Errorf("error create OrganizationMenuTemplate: %v", err)
				return fmt.Errorf("create OrganizationMenuTemplate fail: %w", err)
			}
		}
	}

	// neu co thi chi lay student cua organization do
	for _, studentID := range studentIDs {
		// kiểm tra student có thuộc organization được quản lý không
		isValid, err := receiver.StudentApplicationRepository.CheckStudentBelongsToOrganizations(tx, studentID, orgIDsManaged)
		if err != nil {
			log.Errorf("Error CheckStudentBelongsToOrganizations: %v", err)
			return errors.New("student does not belong to any organization")
		}
		if !isValid {
			continue // bỏ qua student này
		}

		existing, err := receiver.StudentMenuRepository.GetByStudentIDAndComponentID(tx, studentID, componentID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("get student menu fail: %w", err)
		}

		if existing != nil {
			// Đã tồn tại → update
			existing.Order = order
			existing.Visible = visible
			existing.IsShow = true
			if err := receiver.StudentMenuRepository.UpdateWithTx(tx, existing); err != nil {
				return fmt.Errorf("update student menu fail: %w", err)
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
				return fmt.Errorf("create student menu fail: %w", err)
			}
		}
	}
	return nil
}

func (receiver *UploadSectionMenuUseCase) createTeacherMenus(ctx *gin.Context, tx *gorm.DB, componentID uuid.UUID, visible bool, order int, teacherIDs []uuid.UUID, sectionID uuid.UUID) error {
	// Lấy user hiện tại và danh sách organization được quản lý
	user, err := receiver.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return err
	}

	if len(user.Organizations) == 0 {
		return errors.New("user does not belong to any organization")
	}

	orgIDsManaged, err := user.GetManagedOrganizationIDs(receiver.UserEntityRepository.GetDB())
	if err != nil {
		return err
	}
	if len(orgIDsManaged) == 0 {
		return errors.New("user does not manage any organization")
	}

	// Tạo hoặc update OrganizationMenuTemplate cho mỗi tổ chức quản lý
	for _, orgID := range orgIDsManaged {
		existingTemplate, err := receiver.OrganizationMenuTemplateRepository.GetByOrgIDComponentIDSectionID(
			tx,
			orgID,
			componentID,
			sectionID,
		)
		if err != nil {
			log.Printf("Error OrganizationMenuTemplate: %v", err)
			return fmt.Errorf("error OrganizationMenuTemplate fail: %w", err)
		}

		if existingTemplate == nil {
			newTemplate := &entity.OrganizationMenuTemplate{
				ID:             uuid.New().String(),
				OrganizationID: orgID,
				ComponentID:    componentID.String(),
				SectionID:      sectionID.String(),
			}
			if err := receiver.OrganizationMenuTemplateRepository.CreateWithTx(tx, newTemplate); err != nil {
				log.Printf("Error create OrganizationMenuTemplate: %v", err)
				return fmt.Errorf("create OrganizationMenuTemplate fail: %w", err)
			}
		}
	}

	// Lặp qua danh sách teacherIDs
	for _, teacherID := range teacherIDs {
		// Kiểm tra giáo viên có thuộc tổ chức được quản lý hay không
		isValid, err := receiver.TeacherApplicationRepository.CheckTeacherBelongsToOrganizations(tx, teacherID, orgIDsManaged)
		if err != nil {
			log.Printf("Error CheckTeacherBelongsToOrganizations: %v", err)
			return errors.New("teacher does not belong to any organization")
		}
		if !isValid {
			continue
		}

		existing, err := receiver.TeacherMenuRepository.GetByTeacherIDAndComponentID(tx, teacherID, componentID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("get teacher menu fail: %w", err)
		}

		if existing != nil {
			// Đã tồn tại → update
			existing.Order = order
			existing.Visible = visible
			existing.IsShow = true
			if err := receiver.TeacherMenuRepository.UpdateWithTx(tx, existing); err != nil {
				return fmt.Errorf("update teacher menu fail: %w", err)
			}
		} else {
			// Không tồn tại → tạo mới
			menu := &entity.TeacherMenu{
				ID:          uuid.New(),
				TeacherID:   teacherID,
				ComponentID: componentID,
				Order:       order,
				IsShow:      true,
				Visible:     visible,
			}
			if err := receiver.TeacherMenuRepository.CreateWithTx(tx, menu); err != nil {
				return fmt.Errorf("create teacher menu fail: %w", err)
			}
		}
	}
	return nil
}

func (receiver *UploadSectionMenuUseCase) createStaffMenus(ctx *gin.Context, tx *gorm.DB, componentID uuid.UUID, visible bool, order int, staffIDs []uuid.UUID, sectionID uuid.UUID) error {
	// Lấy user hiện tại và danh sách organization được quản lý
	user, err := receiver.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return err
	}

	if len(user.Organizations) == 0 {
		return errors.New("user does not belong to any organization")
	}

	orgIDsManaged, err := user.GetManagedOrganizationIDs(receiver.UserEntityRepository.GetDB())
	if err != nil {
		return err
	}
	if len(orgIDsManaged) == 0 {
		return errors.New("user does not manage any organization")
	}

	// Tạo hoặc update OrganizationMenuTemplate cho mỗi tổ chức quản lý
	for _, orgID := range orgIDsManaged {
		existingTemplate, err := receiver.OrganizationMenuTemplateRepository.GetByOrgIDComponentIDSectionID(
			tx,
			orgID,
			componentID,
			sectionID,
		)
		if err != nil {
			log.Printf("Lỗi khi kiểm tra OrganizationMenuTemplate: %v", err)
			return fmt.Errorf("kiểm tra OrganizationMenuTemplate thất bại: %w", err)
		}

		if existingTemplate == nil {
			newTemplate := &entity.OrganizationMenuTemplate{
				ID:             uuid.New().String(),
				OrganizationID: orgID,
				ComponentID:    componentID.String(),
				SectionID:      sectionID.String(),
			}
			if err := receiver.OrganizationMenuTemplateRepository.CreateWithTx(tx, newTemplate); err != nil {
				log.Printf("Lỗi khi tạo OrganizationMenuTemplate: %v", err)
				return fmt.Errorf("tạo OrganizationMenuTemplate thất bại: %w", err)
			}
		}
	}

	// Lặp qua danh sách teacherIDs
	for _, staffID := range staffIDs {
		// Kiểm tra giáo viên có thuộc tổ chức được quản lý hay không
		isValid, err := receiver.StaffApplicationRepository.CheckStaffBelongsToOrganizations(tx, staffID, orgIDsManaged)
		if err != nil {
			log.Printf("Error CheckTeacherBelongsToOrganizations: %v", err)
			return errors.New("teacher does not belong to any organization")
		}
		if !isValid {
			continue
		}

		existing, err := receiver.StaffMenuRepository.GetByStaffIDAndComponentID(tx, staffID, componentID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("get teacher menu fail: %w", err)
		}

		if existing != nil {
			// Đã tồn tại → update
			existing.Order = order
			existing.Visible = visible
			existing.IsShow = true
			if err := receiver.StaffMenuRepository.UpdateWithTx(tx, existing); err != nil {
				return fmt.Errorf("update teacher menu fail: %w", err)
			}
		} else {
			// Không tồn tại → tạo mới
			menu := &entity.StaffMenu{
				ID:          uuid.New(),
				StaffID:     staffID,
				ComponentID: componentID,
				Order:       order,
				IsShow:      true,
				Visible:     visible,
			}
			if err := receiver.StaffMenuRepository.CreateWithTx(tx, menu); err != nil {
				return fmt.Errorf("create staff menu fail: %w", err)
			}
		}
	}
	return nil
}

func (receiver *UploadSectionMenuUseCase) UploadSectionMenuV2(ctx *gin.Context, req request.UploadSectionMenuRequest) error {
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

	teacherIDs, err := receiver.TeacherApplicationRepository.GetAllTeacherIDs()
	if err != nil {
		logrus.Error("Rollback by error getting teacher_ids:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("Get list teacher_id failed: %w", err)
	}

	staffIDs, err := receiver.StaffApplicationRepository.GetAllStaffIDs()
	if err != nil {
		logrus.Error("Rollback by error getting staff_ids:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("Get list staff_id failed: %w", err)
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
				if err := receiver.createStudentMenus(ctx, tx, componentID, visible, idx, studentIDs, roleOrg.ID); err != nil {
					logrus.Error("rollback by error create student menu:", err)
					tx.Rollback()
					rolledBack = true
					return err
				}
			case string(value.RoleTeacher):
				if err := receiver.createTeacherMenus(ctx, tx, componentID, visible, idx, teacherIDs, roleOrg.ID); err != nil {
					logrus.Error("rollback by error create teacher menu:", err)
					tx.Rollback()
					rolledBack = true
					return err
				}
			case string(value.RoleStaff):
				if err := receiver.createStaffMenus(ctx, tx, componentID, visible, idx, staffIDs, roleOrg.ID); err != nil {
					logrus.Error("rollback by error create staff menu:", err)
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

	// Xóa teacher menu
	if err := receiver.TeacherMenuRepository.DeleteByComponentID(componentID); err != nil {
		return fmt.Errorf("UploadSectionMenuUseCase.DeleteSectionMenu: delete teacher menu failed: %w", err)
	}

	// Xóa staff menu
	if err := receiver.StaffMenuRepository.DeleteByComponentID(componentID); err != nil {
		return fmt.Errorf("UploadSectionMenuUseCase.DeleteSectionMenu: delete staff menu failed: %w", err)
	}

	// Xóa OrganizationMenuTemplate
	if err := receiver.OrganizationMenuTemplateRepository.DeleteByComponentID(componentID); err != nil {
		return fmt.Errorf("UploadSectionMenuUseCase.DeleteSectionMenu: delete organization menu template failed: %w", err)
	}

	return nil
}
