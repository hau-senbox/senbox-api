package usecase

import (
	"fmt"
	"log"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"

	"github.com/google/uuid"
)

type CreateUserFormApplicationUseCase struct {
	*repository.UserEntityRepository
	*repository.RoleOrgSignUpRepository
	*repository.ComponentRepository
	*repository.StudentMenuRepository
	*repository.TeacherMenuRepository
	*repository.StaffMenuRepository
	*repository.OrganizationMenuTemplateRepository
}

func (receiver *CreateUserFormApplicationUseCase) CreateTeacherFormApplication(req request.CreateTeacherFormApplicationRequest) error {
	_, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: req.UserID})
	if err != nil {
		return err
	}

	teacherID := uuid.New()

	err = receiver.UserEntityRepository.CreateTeacherFormApplication(&entity.STeacherFormApplication{
		ID:             teacherID,
		UserID:         uuid.MustParse(req.UserID),
		OrganizationID: uuid.MustParse(req.OrganizationID),
	})

	if err == nil {
		// Lấy role "teacher"
		roleOrgTeacher, _ := receiver.RoleOrgSignUpRepository.GetByRoleName(string(value.RoleTeacher))
		if roleOrgTeacher == nil {
			return nil // Không có role teacher, không cần tạo menu
		}

		sectionTeacherID := roleOrgTeacher.ID
		organizationID := uuid.MustParse(req.OrganizationID)

		// Lấy các Component ID từ bảng OrganizationMenuTemplate theo sectionID và organizationID
		menuTemplates, err := receiver.OrganizationMenuTemplateRepository.GetBySectionIDAndOrganizationID(sectionTeacherID.String(), organizationID.String())
		if err != nil {
			return fmt.Errorf("error get OrganizationMenuTemplate teacher: %w", err)
		}

		for index, template := range menuTemplates {

			// Lấy thông tin component
			component, err := receiver.ComponentRepository.GetByID(template.ComponentID)
			if err != nil {
				log.Printf("Not found component %v: %v", template.ComponentID, err)
				continue
			}

			visible, _ := helper.GetVisibleToValueComponent(string(component.Value))

			// → Tạo mới một Component từ thông tin đã lấy
			newComponent := &components.Component{
				ID:        uuid.New(),
				Name:      component.Name,
				Type:      component.Type,
				Key:       component.Key,
				SectionID: component.SectionID,
				Value:     component.Value,
			}

			err = receiver.ComponentRepository.Create(newComponent)
			if err != nil {
				log.Printf("Create new component fail: %v", err)
				continue
			}

			err = receiver.TeacherMenuRepository.Create(&entity.TeacherMenu{
				ID:          uuid.New(),
				TeacherID:   teacherID,
				ComponentID: newComponent.ID,
				Order:       index,
				IsShow:      true,
				Visible:     visible,
			})

			if err != nil {
				log.Printf("Create TeacherMenu fail %v: %v", newComponent.ID.String(), err)
				continue
			}
		}
	}

	return err
}

func (receiver *CreateUserFormApplicationUseCase) CreateStaffFormApplication(req request.CreateStaffFormApplicationRequest) error {
	_, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: req.UserID})
	if err != nil {
		return err
	}

	staffID := uuid.New()

	err = receiver.UserEntityRepository.CreateStaffFormApplication(&entity.SStaffFormApplication{
		ID:             staffID,
		UserID:         uuid.MustParse(req.UserID),
		OrganizationID: uuid.MustParse(req.OrganizationID),
	})

	if err == nil {
		// Lấy role "staff"
		roleOrgStaff, _ := receiver.RoleOrgSignUpRepository.GetByRoleName(string(value.RoleStaff))
		if roleOrgStaff == nil {
			return nil // Không có role staff, không cần tạo menu
		}

		sectionStaffID := roleOrgStaff.ID
		organizationID := uuid.MustParse(req.OrganizationID)

		// Lấy các Component ID từ bảng OrganizationMenuTemplate theo sectionID và organizationID
		menuTemplates, err := receiver.OrganizationMenuTemplateRepository.GetBySectionIDAndOrganizationID(sectionStaffID.String(), organizationID.String())
		if err != nil {
			return fmt.Errorf("error get OrganizationMenuTemplate staff: %w", err)
		}

		for index, template := range menuTemplates {
			componentID := template.ComponentID

			// Lấy thông tin component
			component, err := receiver.ComponentRepository.GetByID(componentID)
			if err != nil {
				log.Printf("WARNING: Không tìm thấy component %v: %v", componentID, err)
				continue
			}

			visible, _ := helper.GetVisibleToValueComponent(string(component.Value))

			// → Tạo mới một Component từ thông tin đã lấy
			newComponent := &components.Component{
				ID:        uuid.New(),
				Name:      component.Name,
				Type:      component.Type,
				Key:       component.Key,
				SectionID: component.SectionID,
				Value:     component.Value,
			}

			err = receiver.ComponentRepository.Create(newComponent)
			if err != nil {
				log.Printf("Create new component fail: %v", err)
				continue
			}

			err = receiver.StaffMenuRepository.Create(&entity.StaffMenu{
				ID:          uuid.New(),
				StaffID:     staffID,
				ComponentID: newComponent.ID,
				Order:       index,
				IsShow:      true,
				Visible:     visible,
			})

			if err != nil {
				log.Printf("WARNING: Create StaffMenu fail %v: %v", componentID, err)
				continue
			}
		}
	}

	return err
}

func (receiver *CreateUserFormApplicationUseCase) CreateStudentFormApplication(req request.CreateStudentFormApplicationRequest) error {
	_, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: req.UserID})
	if err != nil {
		return err
	}

	studentID := uuid.New()

	err = receiver.UserEntityRepository.CreateStudentFormApplication(&entity.SStudentFormApplication{
		ID:             studentID,
		StudentName:    req.StudentName,
		ChildID:        uuid.MustParse(req.ChildID),
		UserID:         uuid.MustParse(req.UserID),
		OrganizationID: uuid.MustParse(req.OrganizationID),
	})

	if err == nil {
		// Lấy role "student"
		roleOrgStudent, _ := receiver.RoleOrgSignUpRepository.GetByRoleName(string(value.RoleStudent))
		if roleOrgStudent == nil {
			return nil // Không có role student, không cần tạo menu
		}

		sectionStudentID := roleOrgStudent.ID
		organizationID := uuid.MustParse(req.OrganizationID)

		// Lấy các Component ID từ bảng OrganizationMenuTemplate theo sectionID và organizationID
		menuTemplates, err := receiver.OrganizationMenuTemplateRepository.GetBySectionIDAndOrganizationID(sectionStudentID.String(), organizationID.String())
		if err != nil {
			return fmt.Errorf("lỗi khi lấy OrganizationMenuTemplate: %w", err)
		}

		for index, template := range menuTemplates {

			// Lấy thông tin component để tính Visible (nếu cần)
			component, err := receiver.ComponentRepository.GetByID(template.ComponentID)
			if err != nil {
				// log.Warnf("Không tìm thấy component %v: %v", componentID, err)
				continue
			}

			visible, _ := helper.GetVisibleToValueComponent(string(component.Value))

			// → Tạo mới một Component từ thông tin đã lấy
			newComponent := &components.Component{
				ID:        uuid.New(),
				Name:      component.Name,
				Type:      component.Type,
				Key:       component.Key,
				SectionID: component.SectionID,
				Value:     component.Value,
			}

			err = receiver.ComponentRepository.Create(newComponent)
			if err != nil {
				log.Printf("Create new component fail: %v", err)
				continue
			}

			err = receiver.StudentMenuRepository.Create(&entity.StudentMenu{
				ID:          uuid.New(),
				StudentID:   studentID,
				ComponentID: newComponent.ID,
				Order:       index,
				IsShow:      true,
				Visible:     visible,
			})

			if err != nil {
				log.Printf("Create StudentMenu fail %v: %v", newComponent.ID.String(), err)
				continue
			}
		}
	}

	return err
}
