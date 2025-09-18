package usecase

import (
	"errors"
	"fmt"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/entity/menu"
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
	*repository.DeviceMenuRepository
	*repository.ParentMenuRepository
	*repository.TeacherMenuOrganizationRepository
	*repository.DepartmentMenuRepository
	*repository.DepartmentMenuOrganizationRepository
	*repository.SuperAdminEmergencyMenuRepository
	*repository.OrganizationEmergencyMenuRepository
}

// func (receiver *UploadSectionMenuUseCase) createChildMenus(tx *gorm.DB, componentID uuid.UUID, visible bool, order int, childIDs []uuid.UUID) error {
// 	for _, childID := range childIDs {
// 		existing, err := receiver.ChildMenuRepository.GetByChildIDAndComponentID(tx, childID, componentID)
// 		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
// 			return fmt.Errorf("get child menu fail: %w", err)
// 		}

// 		if existing != nil {
// 			// Đã tồn tại: cập nhật
// 			existing.Order = order
// 			existing.Visible = visible
// 			existing.IsShow = true
// 			if err := receiver.ChildMenuRepository.UpdateWithTx(tx, existing); err != nil {
// 				return fmt.Errorf("update child menu fail: %w", err)
// 			}
// 		} else {
// 			// Không tồn tại: tạo mới
// 			menu := &entity.ChildMenu{
// 				ID:          uuid.New(),
// 				ChildID:     childID,
// 				ComponentID: componentID,
// 				Order:       order,
// 				IsShow:      true,
// 				Visible:     visible,
// 			}
// 			if err := receiver.ChildMenuRepository.CreateWithTx(tx, menu); err != nil {
// 				return fmt.Errorf("create child menu fail: %w", err)
// 			}
// 		}
// 	}
// 	return nil
// }

func (receiver *UploadSectionMenuUseCase) createStudentsMenusTemplate(ctx *gin.Context, tx *gorm.DB, componentID uuid.UUID, sectionID uuid.UUID) error {
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

	return nil
}

func (receiver *UploadSectionMenuUseCase) createTeachersMenusTemplate(ctx *gin.Context, tx *gorm.DB, componentID uuid.UUID, sectionID uuid.UUID) error {
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

	return nil
}

func (receiver *UploadSectionMenuUseCase) createStaffsMenusTemplate(ctx *gin.Context, tx *gorm.DB, componentID uuid.UUID, sectionID uuid.UUID) error {
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

		for _, compReq := range item.Components {
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
					// component.ID = uuid.New()
					// componentID = component.ID
					// if err := receiver.ComponentRepository.CreateWithTx(tx, component); err != nil {
					// 	logrus.Error("rollback by error create component (from non-existent id):", err)
					// 	tx.Rollback()
					// 	rolledBack = true
					// 	return fmt.Errorf("create component fail: %w", err)
					// }
					logrus.Error("rollback by error create component (from non-existent id):", err)
					tx.Rollback()
					rolledBack = true
					return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

			switch roleOrg.RoleName {
			// case string(value.RoleChild):
			// 	if err := receiver.createChildMenus(tx, componentID, visible, idx, childIDs); err != nil {
			// 		logrus.Error("rollback by error create child menu:", err)
			// 		tx.Rollback()
			// 		rolledBack = true
			// 		return err
			// 	}
			case string(value.RoleStudent):
				if err := receiver.createStudentsMenusTemplate(ctx, tx, componentID, roleOrg.ID); err != nil {
					logrus.Error("rollback by error create student menu:", err)
					tx.Rollback()
					rolledBack = true
					return err
				}
			case string(value.RoleTeacher):
				if err := receiver.createTeachersMenusTemplate(ctx, tx, componentID, roleOrg.ID); err != nil {
					logrus.Error("rollback by error create teacher menu:", err)
					tx.Rollback()
					rolledBack = true
					return err
				}
			case string(value.RoleStaff):
				if err := receiver.createStaffsMenusTemplate(ctx, tx, componentID, roleOrg.ID); err != nil {
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
	// Xóa organization device menu
	if err := receiver.MenuRepository.DeleteDeviceMenuOrganizationByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa component
	if err := receiver.ComponentRepository.DeleteComponent(componentID, nil); err != nil {
		return nil
	}

	// Xóa child menu
	if err := receiver.ChildMenuRepository.DeleteByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa student menu
	if err := receiver.StudentMenuRepository.DeleteByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa teacher menu
	if err := receiver.TeacherMenuRepository.DeleteByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa staff menu
	if err := receiver.StaffMenuRepository.DeleteByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa parent menu
	if err := receiver.ParentMenuRepository.DeleteByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa device menu
	if err := receiver.DeviceMenuRepository.DeleteByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa teacher menu organization
	if err := receiver.TeacherMenuOrganizationRepository.DeleteByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa super admin menu
	if err := receiver.MenuRepository.DeleteSuperAdminMenuByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa Organization admin menu
	if err := receiver.MenuRepository.DeleteOrganizationAdminMenuByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa OrganizationMenuTemplate
	if err := receiver.OrganizationMenuTemplateRepository.DeleteByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa Department menu
	if err := receiver.DepartmentMenuRepository.DeleteByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa Department menu organization
	if err := receiver.DepartmentMenuOrganizationRepository.DeleteByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa super admin emergency menu
	if err := receiver.SuperAdminEmergencyMenuRepository.DeleteByComponentID(componentID); err != nil {
		return nil
	}

	// Xóa organization emergency menu
	if err := receiver.OrganizationEmergencyMenuRepository.DeleteByComponentID(componentID); err != nil {
		return nil
	}

	return nil
}

func (receiver *UploadSectionMenuUseCase) UploadStudentMenu(ctx *gin.Context, req request.UploadSectionMenuStudentRequest) error {
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

	// 2. Upsert component và tạo menu theo role
	// dau tien xoa component, child menu, student menu neu co mang delete_component_ids
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	roleOrgStudent, err := receiver.RoleOrgSignUpRepository.GetByRoleName(string(value.RoleStudent))
	if err != nil {
		logrus.Error("Rollback by error getting role by role name:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("Get role by role name failed: %w", err)
	}

	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:      compReq.Name,
			Type:      components.ComponentType(compReq.Type),
			Key:       compReq.Key,
			Value:     datatypes.JSON([]byte(compReq.Value)),
			SectionID: roleOrgStudent.ID.String(),
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
				// component.ID = uuid.New()
				// componentID = component.ID
				// if err := receiver.ComponentRepository.CreateWithTx(tx, component); err != nil {
				// 	logrus.Error("rollback by error create component (from non-existent id):", err)
				// 	tx.Rollback()
				// 	rolledBack = true
				// 	return fmt.Errorf("create component fail: %w", err)
				// }
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

		if err := receiver.createStudentMenu(tx, componentID, visible, compReq.Order, uuid.MustParse(req.StudentID), compReq.IsShow); err != nil {
			logrus.Error("rollback by error create student menu:", err)
			tx.Rollback()
			rolledBack = true
			return err
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

func (receiver *UploadSectionMenuUseCase) createStudentMenu(tx *gorm.DB, componentID uuid.UUID, visible bool, order int, studentID uuid.UUID, isShow bool) error {

	// neu co thi chi lay student cua organization do
	existing, err := receiver.StudentMenuRepository.GetByStudentIDAndComponentID(tx, studentID, componentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get student menu fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update
		existing.Order = order
		existing.Visible = visible
		existing.IsShow = isShow
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
			IsShow:      isShow,
			Visible:     visible,
		}
		if err := receiver.StudentMenuRepository.CreateWithTx(tx, menu); err != nil {
			return fmt.Errorf("create student menu fail: %w", err)
		}
	}
	return nil
}

func (receiver *UploadSectionMenuUseCase) UploadTeacherMenu(ctx *gin.Context, req request.UploadSectionMenuTeacherRequest) error {
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

	// 2. Upsert component và tạo menu theo role
	// dau tien xoa component, child menu, student menu neu co mang delete_component_ids
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	roleOrgTeacher, err := receiver.RoleOrgSignUpRepository.GetByRoleName(string(value.RoleTeacher))
	if err != nil {
		logrus.Error("Rollback by error getting role by role name:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("Get role by role name failed: %w", err)
	}

	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:      compReq.Name,
			Type:      components.ComponentType(compReq.Type),
			Key:       compReq.Key,
			Value:     datatypes.JSON([]byte(compReq.Value)),
			SectionID: roleOrgTeacher.ID.String(),
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
				// component.ID = uuid.New()
				// componentID = component.ID
				// if err := receiver.ComponentRepository.CreateWithTx(tx, component); err != nil {
				// 	logrus.Error("rollback by error create component (from non-existent id):", err)
				// 	tx.Rollback()
				// 	rolledBack = true
				// 	return fmt.Errorf("create component fail: %w", err)
				// }
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

		if err := receiver.createTeacherMenu(tx, componentID, visible, compReq.Order, uuid.MustParse(req.TeacherID), compReq.IsShow); err != nil {
			logrus.Error("rollback by error create student menu:", err)
			tx.Rollback()
			rolledBack = true
			return err
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

func (receiver *UploadSectionMenuUseCase) createTeacherMenu(tx *gorm.DB, componentID uuid.UUID, visible bool, order int, teacherID uuid.UUID, isShow bool) error {

	existing, err := receiver.TeacherMenuRepository.GetByTeacherIDAndComponentID(tx, teacherID, componentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get teacher menu fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update
		existing.Order = order
		existing.Visible = visible
		existing.IsShow = isShow
		if err := receiver.TeacherMenuRepository.UpdateWithTx(tx, existing); err != nil {
			return fmt.Errorf("update teacher menu fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		menu := &entity.TeacherMenu{
			ID:          uuid.New(),
			TeacherID:   teacherID,
			ComponentID: componentID,
			Order:       order,
			IsShow:      isShow,
			Visible:     visible,
		}
		if err := receiver.TeacherMenuRepository.CreateWithTx(tx, menu); err != nil {
			return fmt.Errorf("create teacher menu fail: %w", err)
		}
	}
	return nil
}

func (receiver *UploadSectionMenuUseCase) UploadStaffMenu(ctx *gin.Context, req request.UploadSectionMenuStaffRequest) error {
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

	// 2. Upsert component và tạo menu theo role
	// dau tien xoa component, child menu, student menu neu co mang delete_component_ids
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	roleOrgStaff, err := receiver.RoleOrgSignUpRepository.GetByRoleName(string(value.RoleStaff))
	if err != nil {
		logrus.Error("Rollback by error getting role by role name:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("Get role by role name failed: %w", err)
	}

	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:      compReq.Name,
			Type:      components.ComponentType(compReq.Type),
			Key:       compReq.Key,
			Value:     datatypes.JSON([]byte(compReq.Value)),
			SectionID: roleOrgStaff.ID.String(),
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
				// component.ID = uuid.New()
				// componentID = component.ID
				// if err := receiver.ComponentRepository.CreateWithTx(tx, component); err != nil {
				// 	logrus.Error("rollback by error create component (from non-existent id):", err)
				// 	tx.Rollback()
				// 	rolledBack = true
				// 	return fmt.Errorf("create component fail: %w", err)
				// }
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

		if err := receiver.createStaffMenu(tx, componentID, visible, compReq.Order, uuid.MustParse(req.StaffID), compReq.IsShow); err != nil {
			logrus.Error("rollback by error create student menu:", err)
			tx.Rollback()
			rolledBack = true
			return err
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

func (receiver *UploadSectionMenuUseCase) createStaffMenu(tx *gorm.DB, componentID uuid.UUID, visible bool, order int, staffID uuid.UUID, isShow bool) error {

	existing, err := receiver.StaffMenuRepository.GetByStaffIDAndComponentID(tx, staffID, componentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get teacher menu fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update
		existing.Order = order
		existing.Visible = visible
		existing.IsShow = isShow
		if err := receiver.StaffMenuRepository.UpdateWithTx(tx, existing); err != nil {
			return fmt.Errorf("update teacher menu fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		menu := &entity.StaffMenu{
			ID:          uuid.New(),
			StaffID:     staffID,
			ComponentID: componentID,
			Order:       order,
			IsShow:      isShow,
			Visible:     visible,
		}
		if err := receiver.StaffMenuRepository.CreateWithTx(tx, menu); err != nil {
			return fmt.Errorf("create teacher menu fail: %w", err)
		}
	}
	return nil
}

func (receiver *UploadSectionMenuUseCase) UploadChildMenu(ctx *gin.Context, req request.UploadSectionMenuChildRequest) error {
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

	// 2. Upsert component và tạo menu theo role
	// dau tien xoa component, child menu, student menu neu co mang delete_component_ids
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:  compReq.Name,
			Type:  components.ComponentType(compReq.Type),
			Key:   compReq.Key,
			Value: datatypes.JSON([]byte(compReq.Value)),
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
				// component.ID = uuid.New()
				// componentID = component.ID
				// if err := receiver.ComponentRepository.CreateWithTx(tx, component); err != nil {
				// 	logrus.Error("rollback by error create component (from non-existent id):", err)
				// 	tx.Rollback()
				// 	rolledBack = true
				// 	return fmt.Errorf("create component fail: %w", err)
				// }
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

		if err := receiver.createChildMenu(tx, componentID, visible, compReq.Order, uuid.MustParse(req.ChildID), compReq.IsShow); err != nil {
			logrus.Error("rollback by error create child menu:", err)
			tx.Rollback()
			rolledBack = true
			return err
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

func (receiver *UploadSectionMenuUseCase) createChildMenu(tx *gorm.DB, componentID uuid.UUID, visible bool, order int, childID uuid.UUID, isShow bool) error {

	existing, err := receiver.ChildMenuRepository.GetByChildIDAndComponentID(tx, childID, componentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get child menu fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update
		existing.Order = order
		existing.Visible = visible
		existing.IsShow = isShow
		if err := receiver.ChildMenuRepository.UpdateWithTx(tx, existing); err != nil {
			return fmt.Errorf("update child menu fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		menu := &entity.ChildMenu{
			ID:          uuid.New(),
			ChildID:     childID,
			ComponentID: componentID,
			Order:       order,
			IsShow:      isShow,
			Visible:     visible,
		}
		if err := receiver.ChildMenuRepository.CreateWithTx(tx, menu); err != nil {
			return fmt.Errorf("create teacher menu fail: %w", err)
		}
	}
	return nil
}

func (receiver *UploadSectionMenuUseCase) UploadDeviceMenu(ctx *gin.Context, req request.UploadSectionMenuDeviceRequest) error {
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

	// 2. Upsert component và tạo menu theo role
	// dau tien xoa component, child menu, student menu neu co mang delete_component_ids
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:  compReq.Name,
			Type:  components.ComponentType(compReq.Type),
			Key:   compReq.Key,
			Value: datatypes.JSON([]byte(compReq.Value)),
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
				// component.ID = uuid.New()
				// componentID = component.ID
				// if err := receiver.ComponentRepository.CreateWithTx(tx, component); err != nil {
				// 	logrus.Error("rollback by error create component (from non-existent id):", err)
				// 	tx.Rollback()
				// 	rolledBack = true
				// 	return fmt.Errorf("create component fail: %w", err)
				// }
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

		if err := receiver.createDeviceMenu(tx, componentID, visible, compReq.Order, req.DeviceID, compReq.IsShow); err != nil {
			logrus.Error("rollback by error create device menu:", err)
			tx.Rollback()
			rolledBack = true
			return err
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

func (receiver *UploadSectionMenuUseCase) createDeviceMenu(tx *gorm.DB, componentID uuid.UUID, visible bool, order int, deviceID string, isShow bool) error {

	existing, err := receiver.DeviceMenuRepository.GetByDeviceIDAndComponentID(tx, deviceID, componentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get teacher menu fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update
		existing.Order = order
		existing.Visible = visible
		existing.IsShow = isShow
		if err := receiver.DeviceMenuRepository.UpdateWithTx(tx, existing); err != nil {
			return fmt.Errorf("update device menu fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		menu := &entity.SDeviceMenuV2{
			ID:          uuid.New(),
			DeviceID:    deviceID,
			ComponentID: componentID,
			Order:       order,
			IsShow:      isShow,
			Visible:     visible,
		}
		if err := receiver.DeviceMenuRepository.CreateWithTx(tx, menu); err != nil {
			return fmt.Errorf("create teacher menu fail: %w", err)
		}
	}
	return nil
}

func (receiver *UploadSectionMenuUseCase) UploadParentMenu(ctx *gin.Context, req request.UploadSectionMenuParentRequest) error {
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

	// 2. Upsert component và tạo menu theo role
	// dau tien xoa component, child menu, student menu neu co mang delete_component_ids
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:  compReq.Name,
			Type:  components.ComponentType(compReq.Type),
			Key:   compReq.Key,
			Value: datatypes.JSON([]byte(compReq.Value)),
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
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

		if err := receiver.createParentMenu(tx, componentID, visible, compReq.Order, uuid.MustParse(req.ParentID), compReq.IsShow); err != nil {
			logrus.Error("rollback by error create parent menu:", err)
			tx.Rollback()
			rolledBack = true
			return err
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

func (receiver *UploadSectionMenuUseCase) createParentMenu(tx *gorm.DB, componentID uuid.UUID, visible bool, order int, parendID uuid.UUID, isShow bool) error {

	existing, err := receiver.ParentMenuRepository.GetByParentIDAndComponentID(tx, parendID, componentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get child menu fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update
		existing.Order = order
		existing.Visible = visible
		existing.IsShow = isShow
		if err := receiver.ParentMenuRepository.UpdateWithTx(tx, existing); err != nil {
			return fmt.Errorf("update child menu fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		menu := &entity.ParentMenu{
			ID:          uuid.New(),
			ParentID:    parendID,
			ComponentID: componentID,
			Order:       order,
			IsShow:      isShow,
			Visible:     visible,
		}
		if err := receiver.ParentMenuRepository.CreateWithTx(tx, menu); err != nil {
			return fmt.Errorf("create teacher menu fail: %w", err)
		}
	}
	return nil
}

// teacher menu by organization
func (receiver *UploadSectionMenuUseCase) UploadTeacherMenuOrganization(ctx *gin.Context, req request.UploadSectionMenuTeacherOrganizationRequest) error {
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

	// 1. Xóa component nếu có trong danh sách delete_component_ids
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	// 2. Lấy role teacher trong org
	roleOrgTeacher, err := receiver.RoleOrgSignUpRepository.GetByRoleName(string(value.RoleTeacher))
	if err != nil {
		logrus.Error("Rollback by error getting role by role name:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("Get role by role name failed: %w", err)
	}

	// 3. Duyệt qua danh sách component
	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:      compReq.Name,
			Type:      components.ComponentType(compReq.Type),
			Key:       compReq.Key,
			Value:     datatypes.JSON([]byte(compReq.Value)),
			SectionID: roleOrgTeacher.ID.String(),
		}

		if compReq.ID != nil && *compReq.ID != uuid.Nil {
			// Update flow
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
				if err := receiver.ComponentRepository.UpdateWithTx(tx, component); err != nil {
					logrus.Error("rollback by error update component:", err)
					tx.Rollback()
					rolledBack = true
					return fmt.Errorf("update component fail: %w", err)
				}
			} else {
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
			}
		} else {
			// Create flow
			component.ID = uuid.New()
			componentID = component.ID
			if err := receiver.ComponentRepository.CreateWithTx(tx, component); err != nil {
				logrus.Error("rollback by error create component:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail: %w", err)
			}
		}

		// 5. Tạo TeacherMenuOrganization
		if err := receiver.createTeacherMenuOrganization(
			tx,
			componentID.String(),
			compReq.Order,
			req.TeacherID,
			req.OrganizationID,
		); err != nil {
			logrus.Error("rollback by error create teacher menu organization:", err)
			tx.Rollback()
			rolledBack = true
			return err
		}
	}

	// 6. Commit transaction
	if err := tx.Commit().Error; err != nil {
		logrus.Error("Error committing transaction:", err)
		rolledBack = true
		return fmt.Errorf("commit transaction failed: %s", err.Error())
	}

	rolledBack = true
	return nil
}

func (receiver *UploadSectionMenuUseCase) createTeacherMenuOrganization(tx *gorm.DB, componentID string, order int, teacherID string, orgID string) error {

	existing, err := receiver.TeacherMenuOrganizationRepository.GetByTeacherOrgAndComponentID(
		tx,
		teacherID,
		orgID,
		componentID,
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get teacher menu organization fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update order
		existing.Order = order
		if err := receiver.TeacherMenuOrganizationRepository.UpdateWithTx(tx, existing); err != nil {
			return fmt.Errorf("update teacher menu organization fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		menuOrg := &entity.TeacherMenuOrganization{
			TeacherID:      teacherID,
			OrganizationID: orgID,
			ComponentID:    componentID,
			Order:          order,
		}
		if err := receiver.TeacherMenuOrganizationRepository.CreateWithTx(tx, menuOrg); err != nil {
			return fmt.Errorf("create teacher menu organization fail: %w", err)
		}
	}

	return nil
}

// super admin menu
func (receiver *UploadSectionMenuUseCase) UploadSuperAdminMenu(ctx *gin.Context, req request.UploadSectionSuperAdminMenuRequest) error {
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

	// 2. Upsert component và tạo menu theo role
	// dau tien xoa component, child menu, student menu neu co mang delete_component_ids
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:  compReq.Name,
			Type:  components.ComponentType(compReq.Type),
			Key:   compReq.Key,
			Value: datatypes.JSON([]byte(compReq.Value)),
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
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

		if err := receiver.createSuperAdminMenu(tx, componentID.String(), compReq.Order, req.Direction); err != nil {
			logrus.Error("rollback by error create super admin menu:", err)
			tx.Rollback()
			rolledBack = true
			return err
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

func (receiver *UploadSectionMenuUseCase) createSuperAdminMenu(tx *gorm.DB, componentID string, order int, direction menu.Direction) error {

	existing, err := receiver.MenuRepository.GetSuperAdminMenuByComponentID(componentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get super admin menu fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update order
		existing.Order = order
		if err := receiver.MenuRepository.UpdateSuperAdminWithTx(tx, existing); err != nil {
			return fmt.Errorf("update super admin menu fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		superAdminMenu := &menu.SuperAdminMenu{
			Direction:   direction,
			ComponentID: uuid.MustParse(componentID),
			Order:       order,
		}
		if err := receiver.MenuRepository.CreateSuperAdminWithTx(tx, superAdminMenu); err != nil {
			return fmt.Errorf("create super admin menu fail: %w", err)
		}
	}

	return nil
}

// organization admin menu
func (receiver *UploadSectionMenuUseCase) UploadOrganizationAdminMenu(ctx *gin.Context, req request.UploadSectionOrganizationAdminMenuRequest) error {
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

	// 2. Upsert component và tạo menu theo role
	// dau tien xoa component, child menu, student menu neu co mang delete_component_ids
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:  compReq.Name,
			Type:  components.ComponentType(compReq.Type),
			Key:   compReq.Key,
			Value: datatypes.JSON([]byte(compReq.Value)),
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
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

		if err := receiver.createOrganizationAdminMenu(tx, componentID.String(), compReq.Order, req.Direction, req.OrganizationID); err != nil {
			logrus.Error("rollback by error create super admin menu:", err)
			tx.Rollback()
			rolledBack = true
			return err
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

func (receiver *UploadSectionMenuUseCase) createOrganizationAdminMenu(tx *gorm.DB, componentID string, order int, direction menu.Direction, organizationID string) error {

	existing, err := receiver.MenuRepository.GetOrganizationAdminMenuByComponentID(componentID, organizationID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get organization admin menu fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update order
		existing.Order = order
		if err := receiver.MenuRepository.UpdateOrganizationAdminWithTx(tx, existing); err != nil {
			return fmt.Errorf("update organization admin menu fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		organizationAdminMenu := &menu.OrgMenu{
			OrganizationID: uuid.MustParse(organizationID),
			Direction:      direction,
			ComponentID:    uuid.MustParse(componentID),
			Order:          order,
		}
		if err := receiver.MenuRepository.CreateOrganizationAdminWithTx(tx, organizationAdminMenu); err != nil {
			return fmt.Errorf("create organization admin menu fail: %w", err)
		}
	}

	return nil
}

// department menu
func (receiver *UploadSectionMenuUseCase) UploadDepartmentMenu(ctx *gin.Context, req request.UploadSectionMenuDepartmentRequest) error {
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

	// 2. Upsert component và tạo menu theo role
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:  compReq.Name,
			Type:  components.ComponentType(compReq.Type),
			Key:   compReq.Key,
			Value: datatypes.JSON([]byte(compReq.Value)),
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
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

		if err := receiver.createDepartmentMenu(tx, componentID.String(), req.DepartmentID, compReq.Order); err != nil {
			logrus.Error("rollback by error create department menu:", err)
			tx.Rollback()
			rolledBack = true
			return err
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

func (receiver *UploadSectionMenuUseCase) createDepartmentMenu(tx *gorm.DB, componentID string, departmentID string, order int) error {

	existing, err := receiver.DepartmentMenuRepository.GetByDepartmentIDAndComponentID(tx, departmentID, componentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get department menu fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update
		existing.Order = order
		if err := receiver.DepartmentMenuRepository.UpdateWithTx(tx, existing); err != nil {
			return fmt.Errorf("update department menu fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		menu := &entity.DepartmentMenu{
			ID:           uuid.New(),
			DepartmentID: departmentID,
			ComponentID:  uuid.MustParse(componentID),
			Order:        order,
		}
		if err := receiver.DepartmentMenuRepository.CreateWithTx(tx, menu); err != nil {
			return fmt.Errorf("create department menu fail: %w", err)
		}
	}
	return nil
}

// department menu organization
func (receiver *UploadSectionMenuUseCase) UploadDepartmentMenuOrganization(ctx *gin.Context, req request.UploadDepartmentMenuOrganizationRequest) error {
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

	// 2. Upsert component và tạo menu theo role
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:  compReq.Name,
			Type:  components.ComponentType(compReq.Type),
			Key:   compReq.Key,
			Value: datatypes.JSON([]byte(compReq.Value)),
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
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

		if err := receiver.createDepartmentMenuOrganization(tx, componentID.String(), req.DepartmentID, req.OrganizationID, compReq.Order); err != nil {
			logrus.Error("rollback by error create department menu:", err)
			tx.Rollback()
			rolledBack = true
			return err
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

func (receiver *UploadSectionMenuUseCase) createDepartmentMenuOrganization(tx *gorm.DB, componentID string, departmentID string, organizationID string, order int) error {

	existing, err := receiver.DepartmentMenuOrganizationRepository.GetByDepartmentOrgAndComponentID(
		tx,
		departmentID,
		organizationID,
		componentID,
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get department menu organization fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update order
		existing.Order = order
		if err := receiver.DepartmentMenuOrganizationRepository.UpdateWithTx(tx, existing); err != nil {
			return fmt.Errorf("update department menu organization fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		menuOrg := &entity.DepartmentMenuOrganization{
			DepartmentID:   departmentID,
			OrganizationID: organizationID,
			ComponentID:    uuid.MustParse(componentID),
			Order:          order,
		}
		if err := receiver.DepartmentMenuOrganizationRepository.CreateWithTx(tx, menuOrg); err != nil {
			return fmt.Errorf("create department menu organization fail: %w", err)
		}
	}

	return nil
}

// superadmin, organization emergency menu
func (receiver *UploadSectionMenuUseCase) UploadEmergencyMenu(ctx *gin.Context, req request.UploadEmergencyMenuRequest) error {
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

	// 2. Upsert component và tạo menu theo role
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	// dau tien kiem tra user dang la quan ly cua organization nao
	user, err := receiver.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return err
	}

	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:  compReq.Name,
			Type:  components.ComponentType(compReq.Type),
			Key:   compReq.Key,
			Value: datatypes.JSON([]byte(compReq.Value)),
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
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

		if user.IsSuperAdmin() {
			if err := receiver.createSuperAdminEmergencyMenu(tx, componentID.String(), compReq.Order); err != nil {
				logrus.Error("rollback by error create superadmin emergency menu:", err)
				tx.Rollback()
				rolledBack = true
				return err
			}
		} else {
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
			if err := receiver.createOrganizationEmergencyMenu(tx, orgIDsManaged[0], componentID.String(), compReq.Order); err != nil {
				logrus.Error("rollback by error create superadmin emergency menu:", err)
				tx.Rollback()
				rolledBack = true
				return err
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

func (receiver *UploadSectionMenuUseCase) createSuperAdminEmergencyMenu(tx *gorm.DB, componentID string, order int) error {

	existing, err := receiver.SuperAdminEmergencyMenuRepository.GetByComponentID(componentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get super admin emergency menu fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update
		existing.Order = order
		if err := receiver.SuperAdminEmergencyMenuRepository.UpdateWithTx(tx, existing); err != nil {
			return fmt.Errorf("update emergency menu fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		menu := &entity.SuperAdminEmergencyMenu{
			ID:          uuid.New(),
			ComponentID: uuid.MustParse(componentID),
			Order:       order,
		}
		if err := receiver.SuperAdminEmergencyMenuRepository.CreateWithTx(tx, menu); err != nil {
			return fmt.Errorf("create emergency menu fail: %w", err)
		}
	}
	return nil
}

func (receiver *UploadSectionMenuUseCase) createOrganizationEmergencyMenu(tx *gorm.DB, organizationID string, componentID string, order int) error {

	existing, err := receiver.OrganizationEmergencyMenuRepository.GetByOrganizationIDAndComponentID(organizationID, componentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get organization emergency menu fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update
		existing.Order = order
		if err := receiver.OrganizationEmergencyMenuRepository.UpdateWithTx(tx, existing); err != nil {
			return fmt.Errorf("update organization emergency menu fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		menu := &entity.OrganizationEmergencyMenu{
			ID:             uuid.New(),
			ComponentID:    uuid.MustParse(componentID),
			OrganizationID: organizationID,
			Order:          order,
		}
		if err := receiver.OrganizationEmergencyMenuRepository.CreateWithTx(tx, menu); err != nil {
			return fmt.Errorf("create emergency menu fail: %w", err)
		}
	}
	return nil
}

// device menu by org
func (receiver *UploadSectionMenuUseCase) UploadOrganizationDeviceMenu(ctx *gin.Context, req request.UploadSectionDeviceMenuOrganizationRequest) error {
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

	// 2. Upsert component và tạo menu theo role
	if len(req.DeleteComponentIDs) > 0 {
		for _, compID := range req.DeleteComponentIDs {
			if err := receiver.DeleteSectionMenu(compID); err != nil {
				logrus.Error("Rollback by error deleting section menu:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("Delete section menu failed: %w", err)
			}
		}
	}

	for _, compReq := range req.Components {
		var componentID uuid.UUID

		component := &components.Component{
			Name:  compReq.Name,
			Type:  components.ComponentType(compReq.Type),
			Key:   compReq.Key,
			Value: datatypes.JSON([]byte(compReq.Value)),
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
				logrus.Error("rollback by error create component (from non-existent id):", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail (Component ID wrong): %w", err)
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

		if err := receiver.createDeviceMenuOrganization(tx, componentID.String(), req.OrganizationID, compReq.Order); err != nil {
			logrus.Error("rollback by error create device menu organization:", err)
			tx.Rollback()
			rolledBack = true
			return err
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

func (receiver *UploadSectionMenuUseCase) createDeviceMenuOrganization(tx *gorm.DB, componentID string, organizationID string, order int) error {

	existing, err := receiver.MenuRepository.GetDeviceMenuOrganization(
		organizationID,
		componentID,
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get device menu organization fail: %w", err)
	}

	if existing != nil {
		// Đã tồn tại → update order
		existing.Order = order
		if err := receiver.MenuRepository.UpdateDeviceMenuOrganizationWithTx(tx, existing); err != nil {
			return fmt.Errorf("update device menu organization fail: %w", err)
		}
	} else {
		// Không tồn tại → create
		menu := &menu.DeviceMenu{
			OrganizationID: uuid.MustParse(organizationID),
			ComponentID:    uuid.MustParse(componentID),
			Order:          order,
		}
		if err := receiver.MenuRepository.CreateDeviceMenuOrganization(tx, menu); err != nil {
			return fmt.Errorf("create device menu organization fail: %w", err)
		}
	}

	return nil
}
